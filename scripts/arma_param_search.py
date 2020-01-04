# pylama:ignore=E702,E501,E402,E731,E0100
import sys
sys.path.insert(0, '..')
from statsmodels.tsa.arima_model import ARMA
from sklearn.metrics import mean_squared_error
from pymongo import MongoClient
import src.utils.constants as constants
from src.utils import running_fn
import numpy as np
import pandas as pd
import argparse
import warnings
import calendar
import datetime
import math
import pprint
"""
This is needed for clear output, conducting a grid search
there is a 100% chance that there will be parameter combinations
that do not converge for extracted training data. Here we just skip
those, looking for the best combination and saving them in a dictionary.
Without this the output of the script is really cluttery for reason mentioned
earlier. Delete it for your own inconvenience.
"""
warnings.filterwarnings("ignore")

appIDs = '347489791,380077084,748139165,331463193,815616092,486784909,713087156,722018081'
SEARCHED_METRICS = ['conf_count', 'oqmean']


def fluctuation_normal_range_search(dataframe, window_size):
    fluctuation_metric = running_fn(dataframe['value'], np.var, window_size)
    return find_signal_normal_range(fluctuation_metric)


def find_signal_normal_range(signal):
    sorted_signal = np.sort(signal)
    ders = np.diff(sorted_signal)
    res = np.zeros((0, 3))
    min_normal_range_size = int(0.3 * signal.shape[0])
    max_normal_range_size = int(0.95 * signal.shape[0])
    for window_size in range(min_normal_range_size, max_normal_range_size):
        for start in range(0, signal.shape[0] - window_size):
            res = np.vstack((res, [
                window_size, start,
                np.mean(ders[start:start + window_size])
            ]))
    res = res[res[:, 2].astype(int) > 0]
    bestres = res[np.argmin(res[:, 2] / res[:, 0]), :]
    normal_range = sorted_signal[int(bestres[1]):int(bestres[1] + bestres[0])]
    return [np.min(normal_range), np.max(normal_range)]


def parameter_search(appIDs, history, username, password, host, window_size,
                     verbose):
    best_params = {}
    thresholds = {}
    for appID in appIDs:
        timeseries = get_data(appID, history, username, password, host,
                              verbose)
        conf_dict = grid_search(timeseries, appID, verbose)
        best_params[int(appID)] = conf_dict
    show_results(best_params, thresholds)


def get_data(appID, history, username, password, host, verbose):
    dataframe = extract_dataframe(appID, history, username, password, host)
    if verbose:
        print(dataframe)
    return dataframe


def extract_dataframe(appID,
                      history,
                      username,
                      password,
                      host,
                      collection='app_stats'):

    start_dt = datetime.datetime.today().replace(
        hour=0, minute=0, second=0,
        microsecond=0) - datetime.timedelta(days=history)
    df_result = pd.DataFrame()
    for i in range(history):
        print(
            'Fetching {} dataframes {:.2f}%'.format(appID, i / history * 100),
            end='\r')
        q = {}
        start = start_dt + datetime.timedelta(days=i)
        end = start_dt + datetime.timedelta(days=i + 1)
        pipeline = set_pipeline(start, end)
        config = {
            'database': appID,
            'collection': collection,
            'aggregate': pipeline
        }
        client_kwargs = get_client_kwargs(username, password, host)

        _client = MongoClient(**client_kwargs)
        database = config.get("database")
        collection = config.get("collection")
        aggregate = config.get("aggregate")
        client_collection = establish_collection_connection(
            _client, database, collection)

        result = client_collection.aggregate(aggregate)
        q.update({'result': pd.DataFrame(list(result))})
        df_result = pd.concat([df_result, q['result'].copy()])

    df_result['oqmean'] = df_result['objectiveQualityV35'] / df_result['count']
    df_result['conf_count'] = df_result['totalSuccessfulConferences'] + \
        df_result['totalFailedConferences'] + \
        df_result['partiallyFailedConferences'] + \
        df_result['totalDroppedConferences']
    result = df_result.reset_index().copy()
    return result


def datetime_to_epoch_milliseconds(aDateTime):
    return int(calendar.timegm(aDateTime.timetuple()) * constants.CONVERT_KILO)


def get_time_bounds(day_interval, end_time):

    time_delta = datetime.timedelta(days=1)
    start_time = end_time - time_delta

    return start_time, end_time


def set_pipeline(
        start_timestamp,
        end_timestamp,
        project_dict={
            "objectiveQualityV35": "$ObjectiveQualityV35.percentile75",
            "count": "$ObjectiveQualityV35.count",
            "totalSuccessfulConferences": "$totalSuccessfulConferences",
            "totalFailedConferences": "$totalFailedConferences",
            "partiallyFailedConferences": "$partiallyFailedConferences",
            "totalDroppedConferences": "$totalDroppedConferences",
        }):
    pipeline = [{
        "$match": {
            "from": start_timestamp,
            "kind": "daily",
            "to": end_timestamp,
        }
    }, {
        "$project": project_dict
    }]
    return pipeline


def get_client_kwargs(username, password, host):
    result = {
        'host': host,
        'username': username,
        'password': password,
        'authSource': 'admin',
        'authMechanism': 'SCRAM-SHA-1',
        'readPreference': 'secondaryPreferred'
    }
    return result


def establish_collection_connection(_client, database, collection):
    client_db = _client[database]
    client_collection = client_db[collection]
    return client_collection


def grid_search(dataframe, appID, verbose):
    config_dict = {}
    train, test = split_train_test(dataframe, SEARCHED_METRICS)
    for metric in SEARCHED_METRICS:
        print('Starting evalutating metric {} for appID {}'.format(
            metric, appID))
        error = 1e100
        error_p_q_map = {}
        for p in range(1, 10):
            for q in range(10):
                try:
                    recent_error = check_model(train[metric], test[metric], p,
                                               q)
                    print(recent_error)
                    if recent_error < error:
                        error = recent_error
                        error_p_q_map[error] = (p, q)
                        print('Best error: {}\nBest parameters {}'.format(
                            error, error_p_q_map[error]))
                except Exception as e:
                    if verbose:
                        print(e)
                    continue
        if error < 1e100:
            print('Saving {} smallest error for {} metric with p {} and q {}'
                  .format(error, metric, *error_p_q_map[error]))
            config_dict[metric] = error_p_q_map[error]
    return config_dict


def split_train_test(dataframe, searched_metrics):
    dataframe_size = dataframe.shape[0]
    eighty_percent = int(0.8 * dataframe_size)
    train = {}
    test = {}
    for metric in searched_metrics:
        tr = dataframe[metric][:eighty_percent].values.astype(np.float)
        te = dataframe[metric][eighty_percent:].values.astype(np.float)
        train[metric] = tr
        test[metric] = te
    return train, test


def check_model(train, test, p, q):
    model = ARMA(train, order=(p, q))
    model_fit = model.fit(disp=0)
    forecast = model_fit.forecast(steps=test.shape[0])
    return compute_error(test, forecast[0])


def compute_error(test, pred):
    error = math.sqrt(mean_squared_error(test, pred))
    return error


def show_results(best_params, thresholds):
    print('Showing best app_IDs in dictionary form')
    print('ARMAPARAMETERS = ')
    pprint.pprint(best_params)


if __name__ == '__main__':
    parser = argparse.ArgumentParser()

    parser.add_argument(
        '--history',
        type=int,
        default=150,
        help='A number of days to look at while searching')
    parser.add_argument(
        '--username',
        type=str,
        required=True,
        help='Username used for accessing database')
    parser.add_argument(
        '--password',
        type=str,
        required=True,
        help='Password for authentication of the username')
    parser.add_argument(
        '--appids',
        type=str,
        default=appIDs,
        help='List of app ids you want to search params for'
        ' should be of the syntax 1000,1001,1002)'
        ' by default, all appIDs are searched for')
    parser.add_argument(
        '--host',
        type=str,
        default='c1-mongo-router-1.mdb.callstats.io:27017',
        help='Host address of the database')
    parser.add_argument(
        '--window_size',
        type=int,
        default=14,
        help='window size for fluctuation metric, preferably multiple of 7')
    parser.add_argument(
        '--verbose',
        type=bool,
        default=False,
        help='Enable additional logging')
    args = parser.parse_args()

    print('Argument values:')
    print('History: {}'.format(args.history))
    print('Username: {}'.format(args.username))
    print('Password: {}'.format(10 * '*'))
    print('Host: {}'.format(args.host))
    print('AppIDs: {}'.format(args.appids))
    print('Window size: {}'.format(args.window_size))
    print('Verbose flag: {}'.format(args.verbose))

    if args.history > 1000:
        parser.error("history cannot be larger than 1000")

    parameter_search(
        args.appids.split(','), args.history, args.username, args.password,
        args.host, args.window_size, args.verbose)
