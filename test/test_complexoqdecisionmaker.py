from src.decisionmakers import ComplexOQDecisionMaker
import datetime
import itertools
import pandas as pd

ONE_VERDICT = list(
    set([item for item in itertools.permutations([1, 0, 0, 0])]))
TWO_VERDICTS = list(
    set([item for item in itertools.permutations([1, 1, 0, 0])]))
THREE_VERDICTS = list(
    set([item for item in itertools.permutations([1, 1, 1, 0])]))
FOUR_VERDICTS = [1, 1, 1, 1]

MAIN_METRIC = 'objective_quality_v35_average'
APPID = 1
MSG_TYPE = 'trend_detection'

MESSAGE_NOT_SENT = {
    'message_details': (None, ),
    'reasons': [],
    'send': False,
    'type': None,
    'direction': None
}

TREND_TYPE = ['monthly', 'immediate']

MOCK_TIME = pd.to_datetime(
    datetime.datetime(2018, 7, 10, 12, 0, 0, tzinfo=datetime.timezone.utc)),

MOCK_TIME_SEQUENCE = [
    pd.to_datetime(
        datetime.datetime(2018, 7, d, 12, 0, 0, tzinfo=datetime.timezone.utc))
    for d in range(1, 4)
]


def create_none_inputs(verdicts):
    return {
        'objective_quality_v35_average_midterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[0],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': (None, ),
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 0,
                }
            }
        },
        'lossEffectMean_midterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[1],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': (None, ),
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 0,
                }
            }
        },
        'throughputEffectMean_midterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[2],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': (None, ),
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 0,
                }
            }
        },
        'delayEffectMean_midterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[3],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': (None, ),
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 0,
                }
            }
        },
        'lossEffectMean_shortterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[1],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': (None, ),
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 0,
                }
            }
        },
        'throughputEffectMean_shortterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[2],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': (None, ),
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 0,
                }
            }
        },
        'delayEffectMean_shortterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[3],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': (None, ),
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 0,
                }
            }
        },
        'objective_quality_v35_average_shortterm': {
            1: {
                'trend_detection': {
                    None: None,
                }
            }
        },
    }


def create_monthly_inputs(verdicts, time=MOCK_TIME):
    return {
        'objective_quality_v35_average_midterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[0],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5],
                    'trend_monthly_decision': 1,
                    'trend_immediate_decision': 0,
                }
            }
        },
        'lossEffectMean_midterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[1],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5],
                    'trend_monthly_decision': 1,
                    'trend_immediate_decision': 0,
                }
            },
        },
        'throughputEffectMean_midterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[2],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5],
                    'trend_monthly_decision': 1,
                    'trend_immediate_decision': 0,
                }
            },
        },
        'delayEffectMean_midterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[3],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5],
                    'trend_monthly_decision': 1,
                    'trend_immediate_decision': 0,
                }
            }
        },
        'objective_quality_v35_average_shortterm': {
            1: {
                'trend_detection': {
                    None: None,
                }
            }
        },
        'lossEffectMean_shortterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[1],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5],
                    'trend_monthly_decision': 1,
                    'trend_immediate_decision': 0,
                }
            },
        },
        'throughputEffectMean_shortterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[2],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5],
                    'trend_monthly_decision': 1,
                    'trend_immediate_decision': 0,
                }
            },
        },
        'delayEffectMean_shortterm': {
            1: {
                'trend_detection': {
                    'time_model': MOCK_TIME,
                    'value_model': verdicts[3],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5],
                    'trend_monthly_decision': 1,
                    'trend_immediate_decision': 0,
                }
            }
        },
    }


def create_immediate_inputs(verdicts, time=MOCK_TIME):
    return {
        'objective_quality_v35_average_midterm': {
            1: {
                'trend_detection': {
                    'time_model': time,
                    'value_model': verdicts[0],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5, 6],
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 1,
                }
            },
        },
        'lossEffectMean_midterm': {
            1: {
                'trend_detection': {
                    'time_model': time,
                    'value_model': verdicts[1],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5, 6],
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 1,
                }
            },
        },
        'throughputEffectMean_midterm': {
            1: {
                'trend_detection': {
                    'time_model': time,
                    'value_model': verdicts[2],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5, 6],
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 1,
                }
            }
        },
        'delayEffectMean_midterm': {
            1: {
                'trend_detection': {
                    'time_model': time,
                    'value_model': verdicts[3],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5, 6],
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 1,
                }
            },
        },
        'objective_quality_v35_average_shortterm': {
            1: {
                'trend_detection': {
                    None: None,
                }
            }
        },
        'lossEffectMean_shortterm': {
            1: {
                'trend_detection': {
                    'time_model': time,
                    'value_model': verdicts[1],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5, 6],
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 1,
                }
            },
        },
        'throughputEffectMean_shortterm': {
            1: {
                'trend_detection': {
                    'time_model': time,
                    'value_model': verdicts[2],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5, 6],
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 1,
                }
            }
        },
        'delayEffectMean_shortterm': {
            1: {
                'trend_detection': {
                    'time_model': time,
                    'value_model': verdicts[3],
                    'app_id_model': 1,
                    'send': False,
                    'message_details': [1, 2, 3, 4, 5, 6],
                    'trend_monthly_decision': 0,
                    'trend_immediate_decision': 1,
                }
            },
        },
    }


creators = list(
    itertools.chain.from_iterable(
        itertools.repeat([create_monthly_inputs, create_immediate_inputs], 2)))


def setup_model():
    model_out = list()
    model = ComplexOQDecisionMaker(app_ids=[1])
    model.output.connect(model_out.append)
    return model, model_out


# NOTE: No message details, No OQ sent
def test_no_verdicts(expected=MESSAGE_NOT_SENT,
                     input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([0, 0, 0, 0])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_only_OQ(expected=MESSAGE_NOT_SENT, input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([1, 0, 0, 0])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_no_OQ_all_effects(expected=MESSAGE_NOT_SENT,
                           input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([0, 1, 1, 1])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_no_OQ_no_loss(expected=MESSAGE_NOT_SENT,
                       input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([0, 0, 1, 1])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_no_OQ_no_delay(expected=MESSAGE_NOT_SENT,
                        input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([0, 1, 1, 0])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_no_OQ_no_throughput(expected=MESSAGE_NOT_SENT,
                             input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([0, 1, 0, 1])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_no_OQ_only_loss(expected=MESSAGE_NOT_SENT,
                         input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([0, 1, 0, 0])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_no_OQ_only_throughput(expected=MESSAGE_NOT_SENT,
                               input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([0, 0, 1, 0])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_no_OQ_only_delay(expected=MESSAGE_NOT_SENT,
                          input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([0, 0, 0, 1])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


# NOTE: Message details are none (no message to be sent), OQ is given
def test_OQ_and_effects_no_mdetails(expected=MESSAGE_NOT_SENT,
                                    input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([1, 1, 1, 1])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_OQ_no_loss_no_mdetails(expected=MESSAGE_NOT_SENT,
                                input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([1, 0, 1, 1])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_OQ_no_delay_no_mdetails(expected=MESSAGE_NOT_SENT,
                                 input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([1, 1, 1, 0])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_OQ_no_throughput_no_mdetails(expected=MESSAGE_NOT_SENT,
                                      input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([1, 1, 0, 1])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_OQ_and_loss_no_mdetails(expected=MESSAGE_NOT_SENT,
                                 input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([1, 1, 0, 0])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_OQ_and_throughput_no_mdetails(expected=MESSAGE_NOT_SENT,
                                       input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator(verdicts=[1, 0, 1, 0])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


def test_OQ_and_delay_no_mdetails(expected=MESSAGE_NOT_SENT,
                                  input_creator=create_none_inputs):
    model, output = setup_model()
    inputs = input_creator([1, 0, 0, 1])
    model.input(inputs)
    assert len(output) == 1
    assert output[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
        MSG_TYPE] == expected


# NOTE: MESSAGE DETAILS SENT, TREND TYPE: Monthly
def test_no_verdicts_monthly(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([0, 0, 0, 0])
        model.input(inputs)
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_only_OQ_monthly(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([1, 0, 0, 0])
        model.input(inputs)
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_no_OQ_all_effects_monthly(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([0, 1, 1, 1])
        model.input(inputs)
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_no_OQ_no_loss_monthly(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([0, 0, 1, 1])
        model.input(inputs)
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_no_OQ_no_delay_monthly(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([0, 1, 1, 0])
        model.input(inputs)
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_no_OQ_no_throughput_monthly(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([0, 1, 0, 1])
        model.input(inputs)
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_no_OQ_only_loss_monthly(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([0, 1, 0, 0])
        model.input(inputs)
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_no_OQ_only_throughput_monthly(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([0, 0, 1, 0])
        model.input(inputs)
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_no_OQ_only_delay_monthly(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([0, 0, 0, 1])
        model.input(inputs)
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


# NOTE: Monthly messages sent (both upwards and downwards,
#       all directions aligned)
def test_OQ_and_effects_sent(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx](
            [direction, direction, direction, direction])
        model.input(inputs)
        assert (set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                    [APPID][MSG_TYPE]['reasons']) == set([
                        'delayEffectMean_midterm', 'lossEffectMean_midterm',
                        'objective_quality_v35_average',
                        'throughputEffectMean_midterm'
                    ]))
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_no_loss_sent(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, 0, direction, direction])
        model.input(inputs)
        assert (set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                    [APPID][MSG_TYPE]['reasons']) == set([
                        'delayEffectMean_midterm',
                        'objective_quality_v35_average',
                        'throughputEffectMean_midterm'
                    ]))
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
    assert len(output) == 4


def test_OQ_no_delay_sent(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, direction, direction, 0])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([
                       'lossEffectMean_midterm',
                       'objective_quality_v35_average',
                       'throughputEffectMean_midterm'
                   ])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
    assert len(output) == 4


def test_OQ_no_throughput_sent(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, direction, 0, direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([
                       'delayEffectMean_midterm', 'lossEffectMean_midterm',
                       'objective_quality_v35_average'
                   ])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
    assert len(output) == 4


def test_OQ_and_loss_sent(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, direction, 0, 0])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(  # noqa
                       ['lossEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
    assert len(output) == 4


def test_OQ_and_throughput_sent(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, 0, direction, 0])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['objective_quality_v35_average',
                        'throughputEffectMean_midterm'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
    assert len(output) == 4


def test_OQ_and_delay_sent(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, 0, 0, direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['delayEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
    assert len(output) == 4


# NOTE: Monthly messages not sent (both upwards and downwards:
#       all directions conflict with Main direction)
def test_OQ_all_effects_opposite(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx](
            [direction, -direction, -direction, -direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_no_loss_rest_opposite(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, 0, -direction, -direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
    assert len(output) == 4


def test_OQ_no_delay_rest_opposite(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, -direction, -direction, 0])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
    assert len(output) == 4


def test_OQ_no_throughput_opposite(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, -direction, 0, -direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
    assert len(output) == 4


def test_OQ_and_loss_opposite(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, -direction, 0, 0])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
    assert len(output) == 4


def test_OQ_and_throughput_opposite(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, 0, -direction, 0])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
    assert len(output) == 4


def test_OQ_and_delay_opposite(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, 0, 0, -direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
    assert len(output) == 4


# NOTE: One support metric conflicts, one is aligned, one is absent


def test_OQ_loss_ok_throughput_conflict(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, direction, -direction, 0])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['lossEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_loss_ok_delay_conflict(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, direction, 0, -direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['lossEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_throughput_ok_loss_conflict(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, -direction, direction, 0])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['objective_quality_v35_average',
                        'throughputEffectMean_midterm'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_delay_ok_throughput_conflict(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, 0, -direction, direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['delayEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_delay_ok_loss_conflict(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, -direction, 0, direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['delayEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_throughput_ok_delay_conflict(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx]([direction, 0, direction, -direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['objective_quality_v35_average',
                        'throughputEffectMean_midterm'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


# NOTE: Two conflicting support metrics, one aligned
def test_OQ_loss_ok_rest_conflict(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx](
            [direction, direction, -direction, -direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['lossEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_throughput_ok_rest_conflict(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx](
            [direction, -direction, direction, -direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['objective_quality_v35_average',
                        'throughputEffectMean_midterm'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_delay_ok_rest_conflict(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx](
            [direction, -direction, -direction, direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['delayEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


# NOTE: Two support metrics aligned with main, one conflicting
def test_OQ_loss_conflict_rest_ok(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx](
            [direction, -direction, direction, direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([
                       'delayEffectMean_midterm',
                       'objective_quality_v35_average',
                       'throughputEffectMean_midterm'
                   ])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_throughput_conflict_rest_ok(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx](
            [direction, direction, -direction, direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([
                       'lossEffectMean_midterm', 'delayEffectMean_midterm',
                       'objective_quality_v35_average'
                   ])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


def test_OQ_delay_conflict_rest_ok(input_creator=creators):
    model, output = setup_model()
    for idx, direction in enumerate([1, -1, 1, -1]):
        inputs = input_creator[idx](
            [direction, direction, direction, -direction])
        model.input(inputs)
        assert set(output[idx]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([
                       'lossEffectMean_midterm',
                       'objective_quality_v35_average',
                       'throughputEffectMean_midterm'
                   ])
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == list(
                itertools.chain.from_iterable(itertools.repeat(TREND_TYPE,
                                                               2)))[idx]
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True
        assert output[idx]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['direction'] == direction
    assert len(output) == 4


# NOTE: Time Series immediate messages
def test_immediate_no_delayed_OQ_message():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[direction, 0, 0, 0], [0, direction, direction, direction]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        [model.input(inputs[idx]) for idx in range(len(inputs))]
        for output in outputs:
            assert output['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
                MSG_TYPE]['reasons'] == []
            assert output['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
                MSG_TYPE]['type'] == 'immediate'
            assert output['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
                MSG_TYPE]['send'] is False


def test_immediate_all_effects_delayed_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, direction, direction, direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([
                       'delayEffectMean_midterm', 'lossEffectMean_midterm',
                       'objective_quality_v35_average',
                       'throughputEffectMean_midterm'
                   ])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_no_loss_delayed_OQ():
    for direction in [1, -1]:
        model, outputs = setup_model()
        verdicts = [[0, 0, direction, direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([
                       'delayEffectMean_midterm',
                       'objective_quality_v35_average',
                       'throughputEffectMean_midterm'
                   ])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_no_throughput_delayed_OQ():
    for direction in [1, -1]:
        model, outputs = setup_model()
        verdicts = [[0, direction, 0, direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([
                       'lossEffectMean_midterm', 'delayEffectMean_midterm',
                       'objective_quality_v35_average'
                   ])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_no_delay_delayed_OQ():
    for direction in [1, -1]:
        model, outputs = setup_model()
        verdicts = [[0, direction, direction, 0], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([
                       'lossEffectMean_midterm',
                       'objective_quality_v35_average',
                       'throughputEffectMean_midterm'
                   ])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_only_delay_delayed_OQ():
    for direction in [1, -1]:
        model, outputs = setup_model()
        verdicts = [[0, 0, 0, direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['delayEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_only_throughput_delayed_OQ():
    for direction in [1, -1]:
        model, outputs = setup_model()
        verdicts = [[0, 0, direction, 0], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['objective_quality_v35_average',
                        'throughputEffectMean_midterm'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_only_loss_delayed_OQ():
    for direction in [1, -1]:
        model, outputs = setup_model()
        verdicts = [[0, direction, 0, 0], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['lossEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


# NOTE: Conflicting effects - no message sent
def test_immediate_all_effects_conflict_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, -direction, -direction, -direction],
                    [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_immediate_delay_throughput_conflict_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, 0, -direction, -direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_immediate_loss_throughput_conflict_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, -direction, -direction, 0], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_immediate_loss_delay_conflict_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, -direction, 0, -direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_immediate_loss_conflicts_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, -direction, 0, 0], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_immediate_throughput_conflicts_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, 0, -direction, 0], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


def test_immediate_delay_conflicts_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, 0, 0, -direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False


# NOTE: One reason good, one reason bad, one reason absent - send message
def test_immediate_loss_ok_throughput_conflicts_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, direction, -direction, 0], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['lossEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_loss_ok_delay_conflicts_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, direction, 0, -direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['lossEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_delay_ok_loss_conflicts_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, -direction, 0, direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['delayEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_delay_ok_throughput_conflicts_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, 0, -direction, direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['delayEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_thoughput_ok_loss_conflicts_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, -direction, direction, 0], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['objective_quality_v35_average',
                        'throughputEffectMean_midterm'])  # noqa
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_thoughput_ok_delay_conflicts_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, 0, direction, -direction], [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['objective_quality_v35_average',
                        'throughputEffectMean_midterm'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True

# NOTE: Two conflict one aligns the direction - message sent


def test_immediate_thoughput_ok_rest_conflict_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, -direction, direction, -direction],
                    [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['objective_quality_v35_average',
                        'throughputEffectMean_midterm'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_delay_ok_rest_conflict_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, -direction, -direction, direction],
                    [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])

        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['delayEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


def test_immediate_loss_ok_rest_conflict_OQ():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, direction, -direction, -direction],
                    [direction, 0, 0, 0]]
        inputs = [
            create_immediate_inputs(verdicts[idx], MOCK_TIME_SEQUENCE[idx])
            for idx in range(len(verdicts))
        ]
        model.input(inputs[0])
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(inputs[1])
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set(
                       ['lossEffectMean_midterm',
                        'objective_quality_v35_average'])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is True


# NOTE: Testing Flushing mechanism
def test_OQDM_flushes_after_monthly():
    for direction in [-1, 1]:
        model, outputs = setup_model()
        verdicts = [[0, direction, 0, 0], [direction, 0, 0, 0],
                    [direction, 0, 0, 0]]
        input_immediate = create_immediate_inputs(verdicts[0],
                                                  MOCK_TIME_SEQUENCE[0])
        input_monthly = create_monthly_inputs(verdicts[1],
                                              MOCK_TIME_SEQUENCE[1])
        input_after_flush = create_immediate_inputs(verdicts[2],
                                                    MOCK_TIME_SEQUENCE[2])
        model.input(input_immediate)

        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['reasons'] == []
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[0]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(input_monthly)
        assert set(outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'monthly'
        assert outputs[1]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
        model.input(input_after_flush)
        assert set(outputs[2]['midterm']['complex_dm_output'][MAIN_METRIC]
                   [APPID][MSG_TYPE]['reasons']) == set([])
        assert outputs[2]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['type'] == 'immediate'
        assert outputs[2]['midterm']['complex_dm_output'][MAIN_METRIC][APPID][
            MSG_TYPE]['send'] is False
