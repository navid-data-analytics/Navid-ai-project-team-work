from src.pipeline import AppPipelineOQ
import pandas as pd


DEFAULT_DT = pd.Timestamp(2017, 1, 1, 12)
appid = 1
config = {
            'detection': {
                'sliding_window_size': 30,
                'confidence_threshold': 0.95
            },
            'prediction': {
                'history_size': 150,
                'arma_order': {appid: (2, 1)},
            },
            'fluctuation': {appid: {'thresholds': (0, 1)}},
        }

input_df = pd.DataFrame({'value': [1], 'time': [DEFAULT_DT]})
apoq_output_components = set(['trend_detection', 'fluctuation_detection'])


def test_apppipelineoq_builds():
    AppPipelineOQ(config, appid, None)


def test_apppipelineoq_runs():
    result = []
    apoq = AppPipelineOQ(config, appid, None)
    apoq.output.connect(result.append)
    apoq.input(input_df)
    assert set(result[0].keys()) == apoq_output_components
