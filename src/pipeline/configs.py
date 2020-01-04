import logging

logger = logging.getLogger('root')

ARMAParameters = {
    'midterm': {
        'conferences_terminated': {
            331463193: (4, 5),
            347489791: (9, 4),
            380077084: (4, 5),
            486784909: (4, 9),
            722018081: (9, 2),
            748139165: (9, 0),
            815616092: (7, 0),
            739234857: (1,
                        0),  # TODO: 13 Nov 2018, missing 91 days for trainng
            928873129: (5, 1),
            943549171: (8, 1),
            234913325: (9, 1),
            602363212: (1,
                        0),  # TODO: 13 Nov 2018, missing 37 days for training
        }
    },
    'shortterm': {
        'conferences_terminated': {
            331463193: (7, 1),
            347489791: (8, 0),
            380077084: (4, 5),
            486784909: (4, 9),
            722018081: (5, 4),
            748139165: (8, 2),
            815616092: (7, 1),
            739234857: (6, 1),
            928873129: (1, 0),
            943549171: (8, 1),
            234913325: (7, 2),
            602363212: (1, 0)
        }
    }
}

FluctuationParameters = {
    'midterm': {
        'conferences_terminated': {
            331463193: {
                'thresholds': (0.0, 0.2)
            },
            347489791: {
                'thresholds': (0.0, 0.2)
            },
            380077084: {
                'thresholds': (0.0, 0.4)
            },
            486784909: {
                'thresholds': (0.0, 0.25)
            },
            722018081: {
                'thresholds': (0.0, 0.9)
            },
            748139165: {
                'thresholds': (0.0, 0.3)
            },
            815616092: {
                'thresholds': (0.0, 0.3)
            },
            739234857: {
                'thresholds': (0.0, 0.35)
            },
            928873129: {
                'thresholds': (0.0, 0.8)
            },
            943549171: {
                'thresholds': (0.0, 0.55)
            },
            234913325: {
                'thresholds': (0.0, 0.15)
            },
            602363212: {
                'thresholds': (0.0, 0.5)
            },
        },
        'objective_quality_v35_average': {
            331463193: {
                'thresholds': (0.0, 0.06)
            },
            347489791: {
                'thresholds': (0.0, 0.036)
            },
            380077084: {
                'thresholds': (0.0, 0.06)
            },
            486784909: {
                'thresholds': (0.0, 0.02)
            },
            722018081: {
                'thresholds': (0.0, 0.16)
            },
            748139165: {
                'thresholds': (0.0, 0.035)
            },
            815616092: {
                'thresholds': (0.0, 0.04)
            },
            739234857: {
                'thresholds': (0.0, 0.5)
            },
            928873129: {
                'thresholds': (0.0, 0.16)
            },
            943549171: {
                'thresholds': (0.0, 0.08)
            },
            234913325: {
                'thresholds': (0.0, 0.045)
            },
            602363212: {
                'thresholds': (0.0, 0.04)
            },
        },
        'rtt_average': {
            331463193: {
                'thresholds': (0.0, 0.3)
            },
            347489791: {
                'thresholds': (0.0, 0.3)
            },
            380077084: {
                'thresholds': (0.0, 0.3)
            },
            486784909: {
                'thresholds': (0.0, 0.125)
            },
            722018081: {
                'thresholds': (0.0, 1)
            },
            748139165: {
                'thresholds': (0.0, 0.5)
            },
            815616092: {
                'thresholds': (0.0, 0.3)
            },
            739234857: {
                'thresholds': (0.0, 1)
            },
            928873129: {
                'thresholds': (0.0, 0.7)
            },
            943549171: {
                'thresholds': (0.0, 1.1)
            },
            234913325: {
                'thresholds': (0.0, 0.6)
            },
            602363212: {
                'thresholds': (0.0, 0.3)
            },
        }
    },
    'shortterm': {
        'conferences_terminated': {
            234913325: {
                'thresholds': (0.0, 1)
            },
            331463193: {
                'thresholds': (0.0, 0.5)
            },
            347489791: {
                'thresholds': (0.0, 0.24)
            },
            380077084: {
                'thresholds': (0.0, 0.8)
            },
            486784909: {
                'thresholds': (0.0, 0.17)
            },
            602363212: {
                'thresholds': (0.0, 6.5)
            },
            722018081: {
                'thresholds': (0.0, 4.25)
            },
            739234857: {
                'thresholds': (0.0, 1)
            },
            748139165: {
                'thresholds': (0.0, 1)
            },
            815616092: {
                'thresholds': (0.0, 0.8)
            },
            928873129: {
                'thresholds': (0.0, 1)
            },
            943549171: {
                'thresholds': (0.0, 3)
            }
        },
        'rtt_average': {
            234913325: {
                'thresholds': (0.0, 0.5)
            },
            331463193: {
                'thresholds': (0.0, 0.225)
            },
            347489791: {
                'thresholds': (0.0, 0.225)
            },
            380077084: {
                'thresholds': (0.0, 0.3)
            },
            486784909: {
                'thresholds': (0.0, 0.075)
            },
            602363212: {
                'thresholds': (0.0, 0.3)
            },
            722018081: {
                'thresholds': (0.0, 0.65)
            },
            739234857: {
                'thresholds': (0.0, 0.55)
            },
            748139165: {
                'thresholds': (0.0, 0.5)
            },
            815616092: {
                'thresholds': (0.0, 0.47)
            },
            928873129: {
                'thresholds': (0.0, 0.71)
            },
            943549171: {
                'thresholds': (0.0, 0.5)
            }
        },
        'objective_quality_v35_average': {
            234913325: {
                'thresholds': (0.0, 0.035)
            },
            331463193: {
                'thresholds': (0.0, 0.0525)
            },
            347489791: {
                'thresholds': (0.0, 0.07)
            },
            380077084: {
                'thresholds': (0.0, 0.035)
            },
            486784909: {
                'thresholds': (0.0, 0.02)
            },
            602363212: {
                'thresholds': (0.0, 0.06)
            },
            722018081: {
                'thresholds': (0.0, 0.17)
            },
            739234857: {
                'thresholds': (0.0, 0.31)
            },
            748139165: {
                'thresholds': (0.0, 0.07)
            },
            815616092: {
                'thresholds': (0.0, 0.55)
            },
            928873129: {
                'thresholds': (0.0, 0.24)
            },
            943549171: {
                'thresholds': (0.0, 0.12)
            }
        }
    }
}


def PipelineConfig(env):
    model_configs = {
        'midterm': {
            'conferences_terminated': {
                'detection': {
                    'sliding_window_size': 30,
                    'confidence_threshold': 0.8
                },
                'prediction': {
                    'history_size':
                    156,
                    'arma_order':
                    ARMAParameters['midterm']['conferences_terminated'],
                    'seasonality_period':
                    7
                },
                'fluctuation':
                FluctuationParameters['midterm']['conferences_terminated'],
            },
            'objective_quality_v35_average': {
                'detection': {
                    'sliding_window_size': 30,
                    'confidence_threshold': 0.95
                },
                'fluctuation':
                FluctuationParameters['midterm']
                ['objective_quality_v35_average'],
            },
            'rtt_average': {
                'detection': {
                    'sliding_window_size': 30,
                    'confidence_threshold': 0.99
                },
                'fluctuation': FluctuationParameters['midterm']['rtt_average'],
            },
            'delayEffectMean': {
                'detection': {
                    'sliding_window_size': 30,
                    'confidence_threshold': 0.8
                },
            },
            'throughputEffectMean': {
                'detection': {
                    'sliding_window_size': 30,
                    'confidence_threshold': 0.8
                },
            },
            'lossEffectMean': {
                'detection': {
                    'sliding_window_size': 30,
                    'confidence_threshold': 0.8
                },
            },
        },
        'shortterm': {
            'conferences_terminated': {
                'detection': {
                    'sliding_window_size': 7,
                    'confidence_threshold': 0.8
                },
                'prediction': {
                    'history_size':
                    156,
                    'arma_order':
                    ARMAParameters['shortterm']['conferences_terminated'],
                    'seasonality_period':
                    7
                },
                'fluctuation':
                FluctuationParameters['shortterm']['conferences_terminated'],
            },
            'objective_quality_v35_average': {
                'detection': {
                    'sliding_window_size': 7,
                    'confidence_threshold': 0.95
                },
                'fluctuation':
                FluctuationParameters['shortterm']
                ['objective_quality_v35_average'],
            },
            'rtt_average': {
                'detection': {
                    'sliding_window_size': 7,
                    'confidence_threshold': 0.99
                },
                'fluctuation':
                FluctuationParameters['shortterm']['rtt_average'],
            },
            'delayEffectMean': {
                'detection': {
                    'sliding_window_size': 7,
                    'confidence_threshold': 0.8
                },
            },
            'throughputEffectMean': {
                'detection': {
                    'sliding_window_size': 7,
                    'confidence_threshold': 0.8
                },
            },
            'lossEffectMean': {
                'detection': {
                    'sliding_window_size': 7,
                    'confidence_threshold': 0.8
                },
            },
        }
    }
    config = {
        'models': model_configs,
        'data_preprocessor': {
            'data_frame_filter': {
                'column': 'appID',
                'operator': 'equal'
            },
            'data_frame_conf_aggregator': {
                'columns': ['ucID', 'confID'],
            },
            'app_ids': env.app_ids,
        },
        'pipeline_params': {
            '$project': {
                '_id': 0,
                'appID': 1,
                'creationTime': 1,
                'numOutboundVideoStreams': 1,
                'ucID': 1,
                'confID': 1,
                'setupStatus': 1,
            }
        },
        'Scheduler_params': {
            'initial_load': 160
        }
    }
    logger.debug('Pipeline passed config: {}'.format(config))
    return config
