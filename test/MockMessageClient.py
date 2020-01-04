import datetime


class MockMessageClient:
    def __init__(self):
        self._initialize_message_dict()

    def _initialize_message_dict(self):
        self.messages = {
            'conferences_terminated': {
                'trend': {
                    'Immediate_down': self.CreateMockTrendImmediate,
                    '15_days_down': self.CreateMockTrendMonthly,
                    '15_days_up': self.CreateMockTrendMonthly,
                    'Immediate_up': self.CreateMockTrendImmediate,
                    'Immediate_up_short': self.CreateMockTrendImmediate,
                    'Immediate_down_short': self.CreateMockTrendImmediate
                         },
                'prediction': {
                    '15_days_down': self.CreateMockPredictionMonthly,
                    '15_days_up': self.CreateMockPredictionMonthly,
                    '7_days_down': self.CreateMockPredictionMonthly,
                    '7_days_up': self.CreateMockPredictionMonthly
                                },
                'fluctuation': {
                   -1: self.CreateMockFluctuationImmediate,
                    0: self.CreateMockFluctuationImmediate,
                    1: self.CreateMockFluctuationImmediate,
                                },
                'fluctuation_short': {
                   -1: self.CreateMockFluctuationImmediate,
                    0: self.CreateMockFluctuationImmediate,
                    1: self.CreateMockFluctuationImmediate,
                                },
                                        },
            'objective_quality_v35_average': {
                # pylama:ignore=E501
                'trend': {
                     1: {  # TREND UP
                        'immediate': {
                            ('lossEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendImmediate,
                            ('objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendImmediate,
                            ('delayEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendImmediate,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendImmediate,
                            ('lossEffectMean',  'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendImmediate,
                            ('delayEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendImmediate,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendImmediate,
                            },
                        'monthly': {
                            ('lossEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendMonthly,
                            ('objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendMonthly,
                            ('delayEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendMonthly,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendMonthly,
                            ('lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendMonthly,
                            ('delayEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendMonthly,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendMonthly,
                            },
                        },
                    -1: {  # TREND DOWN
                        'immediate': {
                            ('lossEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendImmediate,
                            ('objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendImmediate,
                            ('delayEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendImmediate,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendImmediate,
                            ('lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendImmediate,
                            ('delayEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendImmediate,
                            ('delayEffectMean', 'lossEffectMean',  'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendImmediate,
                            },
                        'monthly': {
                            ('lossEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendMonthly,
                            ('objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendMonthly,
                            ('delayEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendMonthly,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average'): self.CreateMockTrendMonthly,
                            ('lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendMonthly,
                            ('delayEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendMonthly,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateMockTrendMonthly,
                            },
                        },
                    },
                'fluctuation': {
                   -1: self.CreateMockFluctuationImmediate,
                    0: self.CreateMockFluctuationImmediate,
                    1: self.CreateMockFluctuationImmediate,
                                },
                'fluctuation_short': {
                   -1: self.CreateMockFluctuationImmediate,
                    0: self.CreateMockFluctuationImmediate,
                    1: self.CreateMockFluctuationImmediate,
                                },
                    },
            'rtt_average':  {
                'trend': {
                    'Immediate_down': self.CreateMockTrendImmediateRTT, #noqa
                    '15_days_down': self.CreateMockTrendMonthlyRTT,
                    '15_days_up': self.CreateMockTrendMonthlyRTT,
                    'Immediate_up': self.CreateMockTrendImmediateRTT,
                    'Immediate_down_short': self.CreateMockTrendImmediateRTT,
                    'Immediate_up_short': self.CreateMockTrendImmediateRTT
                         },
                'fluctuation': {
                   -1: self.CreateMockFluctuationImmediate,
                    0: self.CreateMockFluctuationImmediate,
                    1: self.CreateMockFluctuationImmediate,
                                },
                'fluctuation_short': {
                   -1: self.CreateMockFluctuationImmediate,
                    0: self.CreateMockFluctuationImmediate,
                    1: self.CreateMockFluctuationImmediate,
                                },
                    }

                }

    def _test_message(self, inputs, expected_types):
        assert len(inputs) == len(expected_types), 'invalid received input length - {} instead of {}'.format(len(inputs), len(expected_types))
        for received, expected in zip(inputs, expected_types):
            assert isinstance(received, expected), 'received input {} of type {} does not match expected types {}'.format(received, type(received), expected)

    # TREND
    def CreateMockTrendMonthly(self, *args):
        expected_types = (datetime.datetime, dict, tuple, int, int)
        self._test_message(args, expected_types)

    def CreateMockTrendImmediate(self, *args):
        expected_types = (datetime.datetime, dict, tuple, int, int, int)
        self._test_message(args, expected_types)

    def CreateMockTrendMonthlyRTT(self, *args):
        expected_types = (datetime.datetime, dict, tuple, int)
        self._test_message(args, expected_types)

    def CreateMockTrendImmediateRTT(self, *args):
        expected_types = (datetime.datetime, dict, tuple, int)
        self._test_message(args, expected_types)

    # FLUCTUATION
    def CreateMockFluctuationImmediate(self, *args):
        expected_types = (datetime.datetime, dict, int)
        self._test_message(args, expected_types)

    # PREDICTIONS
    def CreateMockPredictionMonthly(self, *args):
        expected_types = (datetime.datetime, dict, int, int)
        self._test_message(args, expected_types)
