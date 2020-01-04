from src.pipeline import AppPipelineVolume
import pandas as pd


DEFAULT_DT = pd.Timestamp(2017, 1, 1, 12)
appid = 1
config = {
            'detection': {
                'sliding_window_size': 30,
                'confidence_threshold': 0.8
            },
            'prediction': {
                'history_size': 150,
                'seasonality_period': 7,
                'arma_order': {appid: (2, 1)},
            },
            'fluctuation': {appid: {'thresholds': (0, 1)}},
        }

input_df = pd.DataFrame({'value': [1], 'time': [DEFAULT_DT]})
apv_output_components = set(['trend_prediction', 'trend_detection',
                             'fluctuation_detection'])


def test_apppipelinevolume_builds():
    AppPipelineVolume(config, appid, None)


def test_apppipelinevolume_runs():
    result = []
    apv = AppPipelineVolume(config, appid, None)
    apv.output.connect(result.append)
    apv.input(input_df)
    assert set(result[0].keys()) == apv_output_components
