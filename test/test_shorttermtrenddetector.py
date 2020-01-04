"""Tests model class. In this file two data sets are generated
 and passed into model.The first one contains values.
  The second one contains only zeros as values. """
from src.models import ShortTrendDetector as Detector
import pandas as pd
from collections import deque

model_configuration = {
    "sliding_window_size": 4,
    'confidence_threshold': 0.8
}


def generate_pandas_dataframe():
    df = pd.DataFrame({'year': [2016, 2016, 2016, 2016, 2016, 2016, 2016,
                                2016, 2016, 2016, 2016, 2016, 2016, 2016,
                                2016, 2016, 2016, 2016, 2016, 2016],
                       'month': [2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
                                 2, 2, 2, 2, 2, 2, 2, 2],
                       'day': [1, 2, 3, 4, 5, 6, 7, 8, 9,
                               10, 11, 12, 13, 14, 15, 16,
                               17, 18, 19, 20],
                       'value': [4, 5, 5, 5, 5, 5, 6, 7, 8,
                                 9, 12, 14, 15, 17, 17, 17,
                                 17, 17, 16, 17]})
    df['time'] = pd.to_datetime(df[['year', 'month', 'day']])
    df = df.drop(['year', 'month', 'day'], axis=1)
    return df


def test_model_stores_pandas_dataframe(
        model_out=[], data_container=deque(
            maxlen=model_configuration['sliding_window_size']*2)):
    model = Detector(model_configuration)
    model.output.connect(model_out.append)
    dataframe = generate_pandas_dataframe()
    for _, row in dataframe.iterrows():
        model.input(pd.DataFrame(row).T)
        data_container.append(int(pd.DataFrame(row).T['value']))

    assert len(model_out) == len(dataframe)
    assert data_container == model.sliding_window


def generate_pandas_zeroValue_dataframe():
    df = pd.DataFrame({'year': [2016, 2016, 2016, 2016, 2016, 2016, 2016,
                                2016, 2016, 2016, 2016, 2016, 2016, 2016,
                                2016, 2016, 2016, 2016, 2016, 2016],
                       'month': [2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
                                 2, 2, 2, 2, 2, 2, 2, 2],
                       'day': [1, 2, 3, 4, 5, 6, 7, 8, 9,
                               10, 11, 12, 13, 14, 15, 16,
                               17, 18, 19, 20],
                       'value': [0, 0, 0, 0, 0, 0, 0, 0, 0,
                                 0, 0, 0, 0, 0, 0, 0,
                                 0, 0, 0, 0]})
    df['time'] = pd.to_datetime(df[['year', 'month', 'day']])
    return df


def test_model_zero_value(model_out=[], data_container=deque(
        maxlen=model_configuration['sliding_window_size']*2)):
    model = Detector(model_configuration)
    model.output.connect(model_out.append)
    dataframe = generate_pandas_zeroValue_dataframe()
    for _, row in dataframe.iterrows():
        model.input(pd.DataFrame(row).T)
        data_container.append(int(pd.DataFrame(row).T['value']))

    assert len(model_out) == len(dataframe)
    assert data_container == model.sliding_window
