from datetime import datetime, timezone, timedelta
import sys
import ast
sys.path.append('service/gen/protos')  # noqa
import ai_decision_service_pb2
import ai_decision_service_pb2_grpc
from collections import defaultdict
from src.Grpc.MessageClient import MessageClient
from src.Grpc.CrsClient import CrsClient
from src.Grpc.conversions import \
    datetimeToGrpctimestamp, \
    grpctimestampToDatetime, \
    dictToGrpcdata, \
    grpcdataToDict

from testfixtures import LogCapture

import logging
logger = logging.getLogger('root')
logger.setLevel('INFO')

TEST_APPID = 123456789
TEST_DT = datetime(2018, 6, 14, tzinfo=timezone.utc)
DEFAULT_DT = datetime(1971, 1, 1, tzinfo=timezone.utc)
SUPPRESSION_DT = TEST_DT + timedelta(days=2)
NON_SUPPRESSION_DT = TEST_DT - timedelta(days=2)
TEST_TS = {
    'previous_period': {
        'start': TEST_DT - timedelta(days=-60),
        'end': TEST_DT - timedelta(days=-30)
    },
    'current_period': {
        'start': TEST_DT - timedelta(days=-30),
        'end': TEST_DT
    },
    'future_period': {
        'start': TEST_DT,
        'end': TEST_DT + timedelta(days=30),
    },
}

MOCK_FLAG = {
    'unsuppress': {
        str(TEST_APPID): {
            'MidtermTrend15daysUp': ["20-01-2019", "21-01-2019", "22-01-2019"]
        }
    },
    str(TEST_APPID + 1): {
        'MidtermTrend15daysUp': ["23-01-2018", "24-01-2018", "25-01-2018"]
    },
    'date': {
        str(TEST_APPID): {
            'MidtermTrend15daysUp': ["20-01-2018"]
        },
    }
}

TEST_SCORES = (0, 1)
TEST_CHANGE = 0.75


class mock_service:
    def __init__(self):
        pass

    # State
    def Save(self):
        pass

    def Get(self):
        pass

    # Message
    def Create(self):
        pass

    def List(self):
        pass

    def Aggregate(self):
        pass


def mock_getService(stub):
    return mock_service()


def mock_send(method, request, name, reliable=False):
    return None, None


def test_grpc_crsclient():
    """ tests if the aggregate request to CRS succeeds """
    client = CrsClient('notexistent:5432')
    client._connection_timeout = 0
    client._max_retries = 1
    client.getService = mock_getService
    client.send = mock_send

    with LogCapture() as logs:
        res = client.Aggregate(
            appID=TEST_APPID,
            from_dt=TEST_DT,
            to_dt=TEST_DT + timedelta(days=1))

    assert res is None
    assert 'No logging captured' in str(logs)


def prepare_message_client(*args, **kwargs):
    client = MessageClient('notexistent:5432', *args, **kwargs)
    client._connection_timeout = 0
    client._max_retries = 1
    client.getService = mock_getService
    client.send = mock_send
    return client


def test_grpc_state():
    """ tests if the messages are accepted by gRPC protocols """
    client = prepare_message_client()

    # STATE
    with LogCapture() as logs:
        client.SaveState(keyword='test', state={'state': 'this'})
    assert 'No logging captured' in str(logs)
    with LogCapture() as logs:
        client.GetState(keyword='test')
    assert 'No logging captured' in str(logs)

    with LogCapture() as logs:
        client.SaveState(
            appID=TEST_APPID, keyword='test', state={'state': 'this'})
    assert 'No logging captured' in str(logs)
    with LogCapture() as logs:
        client.GetState(appID=TEST_APPID, keyword='test')
    assert 'No logging captured' in str(logs)

    with LogCapture() as logs:
        client.SaveState(dt=TEST_DT, keyword='test', state={'state': 'this'})
    assert 'No logging captured' in str(logs)
    with LogCapture() as logs:
        client.GetState(dt=TEST_DT, keyword='test')
    assert 'No logging captured' in str(logs)


def test_grpc_create_message():
    """ tests if the messages are accepted by gRPC protocols """
    client = prepare_message_client()

    # TREND
    with LogCapture() as logs:
        client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                                TEST_APPID, 1)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateVolumeMidtermTrend15daysDown(TEST_DT, TEST_TS,
                                                  TEST_SCORES, TEST_APPID, 1)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateVolumeMidtermTrendImmediatelyUp(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, 1, 1)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateVolumeMidtermTrendImmediatelyDown(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, 1, 1)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysUpLoss(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysDownLoss(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysUpDelay(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysDownDelay(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysUpThroughput(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysDownThroughput(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)
    #
    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysUpLossDelay(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysDownLossDelay(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysUpLossThroughput(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysDownLossThroughput(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysUpThroughputDelay(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysDownThroughputDelay(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysUpAll(TEST_DT, TEST_TS, TEST_SCORES,
                                                 TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQCNMidtermTrend15daysDownAll(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE)
    assert 'send' in str(logs)

    # FLUCTUATION
    with LogCapture() as logs:
        client.CreateVolumeMidtermFluctuationImmediatelyHigh(
            TEST_DT, TEST_TS, TEST_APPID)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateVolumeMidtermFluctuationImmediatelyStabilized(
            TEST_DT, TEST_TS, TEST_APPID)
    assert 'send' in str(logs)

    # PREDICTION
    with LogCapture() as logs:
        client.CreateVolumeMidtermPrediction15daysUp(TEST_DT, TEST_TS,
                                                     TEST_APPID, 1)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateVolumeMidtermPrediction15daysDown(TEST_DT, TEST_TS,
                                                       TEST_APPID, 1)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQMidtermTrend15daysDown(TEST_DT, TEST_TS, TEST_SCORES,
                                              TEST_APPID, 1)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                            TEST_APPID, 1)
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQMidtermFluctuationImmediatelyHigh(
            TEST_DT,
            TEST_TS,
            TEST_APPID,
        )
    assert 'send' in str(logs)

    with LogCapture() as logs:
        client.CreateOQMidtermFluctuationImmediatelyStabilized(
            TEST_DT,
            TEST_TS,
            TEST_APPID,
        )
    assert 'send' in str(logs)


def mock_send_list(method, request, name, reliable=False):
    msg1 = ai_decision_service_pb2.Message(
        message='this is msg1',
        app_id=TEST_APPID,
        type='test',
        version=1,
        data=dictToGrpcdata({
            'testdata': 1
        }),
        generation_time=datetimeToGrpctimestamp(TEST_DT),
    )

    msg2 = ai_decision_service_pb2.Message(
        message='this is msg2',
        app_id=TEST_APPID,
        type='test',
        version=1,
        data=dictToGrpcdata({
            'testdata': 2
        }),
        generation_time=datetimeToGrpctimestamp(TEST_DT),
    )

    return [msg1, msg2], None


def test_grpc_list_messages():
    """ tests if the messages are accepted by gRPC protocols """
    client = MessageClient('notexistent:5432')
    client._connection_timeout = 0
    client._max_retries = 1
    client.getService = mock_getService
    client.send = mock_send_list

    with LogCapture() as logs:
        it = client.ListMessages(TEST_APPID)
        for e in it:
            logger.info(e['message'])

    assert 'this is msg1' in str(logs)
    assert 'this is msg2' in str(logs)


def test_message_suppression():
    client = prepare_message_client()
    client._latest_dates = defaultdict(
        dict, {str(TEST_APPID): {
                   "MidtermTrend15daysUp": SUPPRESSION_DT
               }})
    with LogCapture() as logs:
        result = client.CreateVolumeMidtermTrend15daysUp(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, 1)
    assert 'not suppressed' not in str(logs)
    assert result is None


def test_no_message_suppression():
    client = prepare_message_client()
    client._latest_dates = defaultdict(
        dict, {str(TEST_APPID): {
                   "MidtermTrend15daysUp": NON_SUPPRESSION_DT
               }})
    with LogCapture() as logs:
        result = client.CreateVolumeMidtermTrend15daysUp(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, 1)
    assert 'not suppressed' in str(logs)
    assert result is None


def test_suppression_state_no_update():
    client = prepare_message_client()
    original_state = {
        str(TEST_APPID): {
            "MidtermTrend15daysUp": SUPPRESSION_DT
        }
    }
    client._latest_dates = defaultdict(dict, original_state)
    client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                            TEST_APPID, 1)
    assert client._latest_dates == original_state


def test_suppression_state_update():
    client = prepare_message_client()
    original_state = {
        str(TEST_APPID): {
            "MidtermTrend15daysUp": NON_SUPPRESSION_DT
        }
    }
    client._latest_dates = defaultdict(dict, original_state)
    client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                            TEST_APPID, 1)
    client.save_dates(TEST_DT)
    assert client._latest_dates[str(
        TEST_APPID)]["MidtermTrend15daysUp"] == TEST_DT


def test_suppression_state_add_new_appid():
    client = prepare_message_client()
    assert str(TEST_APPID) not in client._latest_dates.keys()
    client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                            TEST_APPID, 1)
    assert str(TEST_APPID) in client._latest_dates.keys()


def test_suppression_state_add_new_type():
    client = prepare_message_client()
    assert "MidtermTrend15daysUp" not in client._latest_dates[str(
        TEST_APPID)].keys()
    client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                            TEST_APPID, 1)
    assert "MidtermTrend15daysUp" in client._latest_dates[str(
        TEST_APPID)].keys()


def test_get_state_valid():
    client = prepare_message_client()
    mock_state = {'state': TEST_DT}
    mock_states = {
        'state': {
            str(TEST_APPID): {
                'MidtermTrend15daysUp': TEST_DT
            }
        }
    }

    def mock_get_state(tag):
        if tag == 'latest_date':
            return mock_state
        elif tag == 'latest_dates':
            return mock_states

    client.GetState = mock_get_state
    result = client._get_dates_from_state()
    assert result == mock_states['state']


def test_get_state_fallback():
    client = prepare_message_client()
    mock_state = {'state': TEST_DT}
    mock_states = None

    def mock_get_state(tag):
        if tag == 'latest_date':
            return mock_state
        elif tag == 'latest_dates':
            return mock_states

    client.GetState = mock_get_state
    result = client._get_dates_from_state()
    assert result == {}
    assert client._default_date == mock_state['state']


def test_get_state_not_available():
    client = prepare_message_client()
    mock_states = None
    mock_state = None

    def mock_get_state(tag):
        if tag == 'latest_date':
            return mock_state
        elif tag == 'latest_dates':
            return mock_states

    client.GetState = mock_get_state
    result = client._get_dates_from_state()
    assert result == {}
    assert client._default_date == DEFAULT_DT


def test_get_state_faulty_entry_deletion():
    client = prepare_message_client()
    state = {TEST_APPID: {'a': TEST_DT, 'b': 'faulty_entry', 'c': TEST_DT}}
    result = client._del_faulty_state_entries(state)
    assert set(result[TEST_APPID].keys()) == set(['a', 'c'])


def test_grpc_conversions():
    dt = datetime.now(timezone.utc)
    grpctimestamp = datetimeToGrpctimestamp(dt)
    assert dt == grpctimestampToDatetime(grpctimestamp)

    dic = {
        "int": 0,
        "float": 1.0,
        "str": "s",
        "None": None,
        "dict": {},
        "dt": dt,
        "date": TEST_DT
    }
    grpcdata = dictToGrpcdata(dic)
    assert dic == grpcdataToDict(grpcdata)


def test_grpc_unreliable():
    client = MessageClient('notexistent:5432')
    client._connection_timeout = 0
    request = ai_decision_service_pb2.StateSaveRequest()

    with LogCapture() as logs:
        service = client.getService(
            ai_decision_service_pb2_grpc.AIDecisionStateServiceStub)
        ret, err = client.send(
            service.Save, request, 'SaveState', reliable=False)

    assert None is ret
    assert 'UNAVAILABLE' in str(logs)
    assert 'retry' not in str(logs)


def test_grpc_reliable():
    client = MessageClient('notexistent:5432')
    client._connection_timeout = 0
    client._max_retries = 1
    request = ai_decision_service_pb2.StateSaveRequest()

    with LogCapture() as logs:
        service = client.getService(
            ai_decision_service_pb2_grpc.AIDecisionStateServiceStub)
        ret, err = client.send(
            service.Save, request, 'SaveState', reliable=True)

    assert None is ret
    assert None is not err
    assert '(UNAVAILABLE) retry in 1s' in str(logs)
    assert '(UNAVAILABLE) too many retries' in str(logs)


def test_grpc_manual_date():
    client = prepare_message_client()
    assert client._latest_dates[str(TEST_APPID)] == {}, client._latest_dates
    client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                            TEST_APPID, 1)
    assert client._latest_dates[str(
        TEST_APPID)]['MidtermTrend15daysUp'] == TEST_DT, client._latest_dates
    client._flags = MOCK_FLAG
    client._handle_manual_update()
    expected_DT = datetime.strptime('20-01-2018',
                                    '%d-%m-%Y').replace(tzinfo=timezone.utc)
    assert client._latest_dates[str(TEST_APPID)][
        'MidtermTrend15daysUp'] == expected_DT, client._latest_dates


def test_grpc_unsuppression_date():
    # 1. Prepare message client
    # 2. Set suppression date to 20-01-2019
    # 3. Try sending message with earlier date
    # 4. Assert suppressed message in logs
    # 5. Unsuppress Test date
    # 6. Send message with test date, assert non suppression
    # 7. Assert suppression date did not change
    initial_flag = {
        'date': {
            str(TEST_APPID): {
                'MidtermTrend15daysUp': ["20-01-2019"]
            },
        },
    }
    client = prepare_message_client()
    client._flags = initial_flag
    client._handle_manual_update()
    assert client._latest_dates[str(
        TEST_APPID)]['MidtermTrend15daysUp'] == datetime.strptime(
            '20-01-2019',
            '%d-%m-%Y').replace(tzinfo=timezone.utc), client._latest_dates
    unsuppress_flag = {
        'unsuppress': {
            str(TEST_APPID): {
                'MidtermTrend15daysUp':
                ["14-06-2018", "15-06-2018", "16-06-2018"]
            }
        },
    }
    with LogCapture() as logs:
        client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                                TEST_APPID, 1)

    assert '123456789 not suppressed' not in str(logs), str(logs)

    client._flags = unsuppress_flag
    client._handle_manual_update()

    with LogCapture() as logs:
        client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                                TEST_APPID, 1)

    assert '123456789 suppressed' not in str(logs), str(logs)
    assert client._latest_dates[str(
        TEST_APPID)]['MidtermTrend15daysUp'] == datetime.strptime(
            '20-01-2019',
            '%d-%m-%Y').replace(tzinfo=timezone.utc), client._latest_dates


def test_grpc_unsuppress_different_types():
    # For different appids
    # Unsuppress for 'MidtermTrend15daysUp' on 20-01-2011
    # Unsuppress for 'MidtermTrend15daysUp' on 23-01-2011
    MOCK_FLAG = {
        'unsuppress': {
            str(TEST_APPID): {
                'MidtermTrend15daysUp': ["20-01-2011"]
            },
            str(TEST_APPID + 1): {
                'MidtermTrend15daysDown': ["23-01-2011"]
            }
        }
    }

    # Set client with suppression date after unsuppress
    client = prepare_message_client()
    client._latest_dates = defaultdict(
        dict, {
            str(TEST_APPID): {
                "MidtermTrend15daysUp": SUPPRESSION_DT
            },
            str(TEST_APPID + 1): {
                "MidtermTrend15daysDown": SUPPRESSION_DT
            }
        })

    # Provide flags to the client
    client._flags = MOCK_FLAG
    client._handle_manual_update()

    # Send messages for 20-01-2011
    TEST_DT = datetime(2011, 1, 20, tzinfo=timezone.utc)
    with LogCapture() as logs:
        client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                                TEST_APPID, 1)
        client.CreateVolumeMidtermTrend15daysDown(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID + 1, 1)

    # Assert suppressed is only one type, another is passed along
    assert '123456789 not suppressed' in str(logs), str(logs)
    assert '123456790 suppressed' in str(logs), str(logs)

    # Send messages for 23-01-2011
    TEST_DT = datetime(2011, 1, 23, tzinfo=timezone.utc)
    with LogCapture() as logs:
        client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                                TEST_APPID, 1)
        client.CreateVolumeMidtermTrend15daysDown(
            TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID + 1, 1)

    # Assert suppressed is only one type, another is passed along
    assert '123456789 suppressed' in str(logs), str(logs)
    assert '123456790 not suppressed' in str(logs), str(logs)


def test_grpc_empty_date_flag():
    client = prepare_message_client()
    assert client._latest_dates[str(TEST_APPID)] == {}, client._latest_dates
    client.CreateVolumeMidtermTrend15daysUp(TEST_DT, TEST_TS, TEST_SCORES,
                                            TEST_APPID, 1)
    flag = {
        'unsuppress': {
            str(TEST_APPID): {
                'MidtermTrend15daysUp':
                ["20-01-2019", "21-01-2019", "22-01-2019"]
            }
        },
        'date': {}
    }
    client._flags = flag
    client._handle_manual_update()


def test_grpc_empty_unsuppression_flag():
    client = prepare_message_client()
    assert client._latest_dates[str(TEST_APPID)] == {}, client._latest_dates
    flag = {
        'unsuppress': {},
        'date': {
            str(TEST_APPID): {
                'MidtermTrend15daysUp':
                ["20-01-2019", "21-01-2019", "22-01-2019"]
            }
        }
    }
    client._flags = flag
    client._handle_manual_update()


def test_grpc_both_flags_empty():
    client = prepare_message_client()
    assert client._latest_dates[str(TEST_APPID)] == {}, client._latest_dates
    flag = {'unsuppress': {}, 'date': {}}
    client._flags = flag
    client._handle_manual_update()


def test_grpc_date_flag_missing():
    client = prepare_message_client()
    assert client._latest_dates[str(TEST_APPID)] == {}, client._latest_dates
    flag = {
        'unsuppress': {
            str(TEST_APPID): {
                'MidtermTrend15daysUp':
                ["20-01-2019", "21-01-2019", "22-01-2019"]
            }
        },
    }
    client._flags = flag
    client._handle_manual_update()


def test_grpc_unsuppression_flag_missing():
    client = prepare_message_client()
    assert client._latest_dates[str(TEST_APPID)] == {}, client._latest_dates
    flag = {
        'date': {
            str(TEST_APPID): {
                'MidtermTrend15daysUp':
                ["20-01-2019", "21-01-2019", "22-01-2019"]
            }
        }
    }
    client._flags = flag
    client._handle_manual_update()


def test_grpc_both_flags_missing():
    client = prepare_message_client()
    assert client._latest_dates[str(TEST_APPID)] == {}, client._latest_dates
    flag = {}
    client._flags = flag
    client._handle_manual_update()


def test_client_flag_propagation():
    unsuppress = '{7233: {\'MidtermTrend15daysDown\': ["20-01-2018"]}}'
    manual_date = '{7233: {\'MidtermTrend15daysUp\': ["20-01-2018"]}}'
    flags = {
        'unsuppress': ast.literal_eval(unsuppress),
        'date': ast.literal_eval(manual_date)
    }
    client = MessageClient('notexistent:5432', flags=flags)
    assert client.flags == flags


def test_client_empty_flag_propagation():
    unsuppress = '{}'
    manual_date = '{}'
    flags = {
        'unsuppress': ast.literal_eval(unsuppress),
        'date': ast.literal_eval(manual_date)
    }
    client = MessageClient('notexistent:5432', flags=flags)
    assert client.flags == flags
