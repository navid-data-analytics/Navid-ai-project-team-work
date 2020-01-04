"""Test ArmaPredictor class."""
from src.models import ShorttermArmaPredictor as ArmaPredictor #noqa
from src.models.ShorttermArmaPredictor import output_prometheus
import pandas as pd
import numpy as np

np.random.seed(0)


def generate_pandas_dataframe():
    df = pd.DataFrame({
        'year': [2016] * 30,
        'month': [5] * 30,
        'day': [d for d in range(1, 31)],
        'value': [
            4, 5, 5, 5, 5, 5, 6, 7, 8, 9, 12, 14, 15, 17, 17, 17, 17, 17, 16,
            17, 18, 18, 18, 17, 19, 19, 18, 19, 20, 0
        ]
    })
    df['time'] = pd.to_datetime(df[['year', 'month', 'day']])
    df = df.drop(['year', 'month', 'day'], axis=1)
    return df


def test_arma_takes_default_order(appID=0,
                                  model_conf={
                                         "history_size": 10,
                                         "seasonality_period": 7,
                                         "arma_order": {
                                             0: (6, 1)
                                         }}):
    arma = ArmaPredictor(model_conf, app_id=None)
    assert arma.order == (1, 0)  # default value if appID not found


def test_arma_takes_specified_order(appID=0,
                                    model_conf={
                                         "history_size": 10,
                                         "seasonality_period": 7,
                                         "arma_order": {
                                             0: (6, 1)
                                         }}):
    arma = ArmaPredictor(model_conf, app_id=appID)
    assert arma.order == model_conf["arma_order"][0]


def test_arma_triggered_pandas_dataframe(appID=0,
                                         model_conf={"history_size": 25,
                                                     "seasonality_period": 7,
                                                     "arma_order": {
                                                        0: (1, 0)}},
                                         model_out=list()):
    arma = ArmaPredictor(model_conf, app_id=appID)
    arma.output.connect(model_out.append)
    dataframe = generate_pandas_dataframe()
    for index, row in dataframe.iterrows():
        arma.input(pd.DataFrame(row).T)
    times_triggered = dataframe.shape[0]
    assert len(model_out) == times_triggered


def test_arma_fills_history_before_forecast(
        appID=0,
        model_conf={"history_size": 25,
                    "seasonality_period": 7,
                    "arma_order": {
                        0: (1, 0)
                    }},
        model_out=list(),
        expected_predictions=6):
    # Generated 30 days, model starts at day 25, expected_predictions=6
    arma = ArmaPredictor(model_conf, app_id=appID)
    arma.output.connect(model_out.append)
    dataframe = generate_pandas_dataframe()
    data_stream = dataframe.shape[0]
    for index, row in dataframe.iterrows():
        arma.input(pd.DataFrame(row).T)
    none_predictions = len([v for v in model_out if v[0] is None])
    num_predictions = data_stream - none_predictions
    assert num_predictions == expected_predictions


def test_arma_not_full_no_prediction(
        appID=0,
        model_conf={"history_size": 31,
                    "seasonality_period": 7,
                    "arma_order": {
                        0: (6, 1)
                    }},
        model_out=list(),
        expected_predictions=0):
    arma = ArmaPredictor(model_conf, app_id=appID)
    arma.output.connect(model_out.append)
    dataframe = generate_pandas_dataframe()
    data_stream = dataframe.shape[0]
    for _, row in dataframe.iterrows():
        arma.input(pd.DataFrame(row).T)
    none_predictions = len([v for v in model_out if v[0] is None])
    num_predictions = data_stream - none_predictions
    assert num_predictions == expected_predictions


def test_arma_handles_zero_on_input(appID=0,
                                    model_conf={"history_size": 25,
                                                "seasonality_period": 7,
                                                "arma_order": {
                                                 0: (3, 0)}},
                                    model_out=list()):
    arma = ArmaPredictor(model_conf, app_id=appID)
    arma.output.connect(model_out.append)
    dataframe = generate_pandas_dataframe()
    for _, row in dataframe.iterrows():
        arma.input(pd.DataFrame(row).T)
    times_triggered = dataframe.shape[0]
    assert len(model_out) == times_triggered
