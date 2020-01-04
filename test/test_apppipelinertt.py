from src.pipeline import AppPipelineRTT
import pandas as pd


DEFAULT_DT = pd.Timestamp(2017, 1, 1, 12)
appid = 1
config = {
            'detection': {
                'sliding_window_size': 30,
                'confidence_threshold': 0.95
            },
            'fluctuation': {appid: {'thresholds': (0, 1)}},
        }

input_df = pd.DataFrame({'value': [1], 'time': [DEFAULT_DT]})
aprtt_output_components = set(['trend_detection',
                               'fluctuation_detection'])


def test_apppipelinertt_builds():
    AppPipelineRTT(config, appid, None)


def test_apppipelinertt_runs():
    result = []
    aprtt = AppPipelineRTT(config, appid, None)
    aprtt.output.connect(result.append)
    aprtt.input(input_df)
    assert set(result[0].keys()) == aprtt_output_components
