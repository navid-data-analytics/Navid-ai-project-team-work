"""Tests notifier class."""
from src.pipeline import Notifier
from src.Grpc.MessageClient import MessageClient
from datetime import datetime, timezone, timedelta
import pandas as pd
import logging

logger = logging.getLogger('root')
logger.setLevel('INFO')
DATE = pd.Timestamp(2018, 2, 10, 1)
TEST_APPID = 427974000
TEST_DT = datetime(2018, 6, 14, tzinfo=timezone.utc)
TEST_SCORES = (0, 1)
TEST_DAYS_TO_REPORT = 60
TEST_CHANGE = .76
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

reasons = {
    'all': [
        'delayEffectMean_midterm', 'lossEffectMean_midterm',
        'objective_quality_v35_average', 'throughputEffectMean_midterm'
    ],
    'loss': ['lossEffectMean_midterm', 'objective_quality_v35_average'],
    'delay': ['delayEffectMean_midterm', 'objective_quality_v35_average'],
    'throughput':
    ['objective_quality_v35_average', 'throughputEffectMean_midterm'],
    'throughput_loss': [
        'lossEffectMean_midterm', 'objective_quality_v35_average',
        'throughputEffectMean_midterm'
    ],
    'throughput_delay': [
        'delayEffectMean_midterm', 'objective_quality_v35_average',
        'throughputEffectMean_midterm'
    ],
    'loss_delay': [
        'delayEffectMean_midterm', 'lossEffectMean_midterm',
        'objective_quality_v35_average'
    ],
}


def mock_message(dt, TEST_APPID, type, version, data):
    return type


UP = 1
DOWN = -1
mock_aid_service = MessageClient(address=None)
# NOTE: Aid service is tested somewhere else
mock_aid_service.save_dates = lambda x: None
mock_aid_service._CreateMessageConsiderState = mock_message


def generate_mock(message_details, trend_type, reasons, send, direction):
    mock_data = {
        'shortterm': {
            'date': TEST_DT,
            'complex_dm_output': {
                'conferences_terminated': {
                    TEST_APPID: {
                        'trend_detection': {
                            'time_model': DATE,
                            'value_model': 0,
                            'app_id_model': TEST_APPID,
                            'app_status': '',
                            'decision': 0
                        }
                    }
                },
                'objective_quality_v35_average': {
                    TEST_APPID: {
                        'trend_detection': {
                            'message_details': message_details,
                            'type': trend_type,
                            'reasons': reasons,
                            'send': send,
                            'direction': direction
                        },
                    }
                }
            }
        },
        'midterm': {
            'date': TEST_DT,
            'complex_dm_output': {
                'conferences_terminated': {
                    TEST_APPID: {
                        'trend_detection': {
                            'time_model': DATE,
                            'value_model': 0,
                            'app_id_model': TEST_APPID,
                            'app_status': '',
                            'decision': 0
                        }
                    }
                },
                'objective_quality_v35_average': {
                    TEST_APPID: {
                        'trend_detection': {
                            'message_details': message_details,
                            'type': trend_type,
                            'reasons': reasons,
                            'send': send,
                            'direction': direction
                        },
                    }
                }
            }
        }
    }
    return mock_data


# def test_notifier_send_oq_up_all_reasons():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['all'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysUpAll'
#
#
# def test_notifier_send_oq_down_all_reasons():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['all'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysDownAll'

#
# def test_notifier_send_oq_up_loss():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['loss'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysUpLoss'
#
#
# def test_notifier_send_oq_down_loss():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['loss'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysDownLoss'
#
#
# def test_notifier_send_oq_up_delay():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['delay'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysUpDelay'
#
#
# def test_notifier_send_oq_down_delay():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['delay'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysDownDelay'
#
#
# def test_notifier_send_oq_up_throughput():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['throughput'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysUpThroughput'
#
#
# def test_notifier_send_oq_down_throughput():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['throughput'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysDownThroughput'
#
#
# def test_notifier_send_oq_down_loss_delay():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['loss_delay'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysDownLossDelay'
#
#
# def test_notifier_send_oq_up_loss_delay():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['loss_delay'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysUpLossDelay'
#
#
# def test_notifier_send_oq_down_throughput_delay():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['throughput_delay'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysDownThroughputDelay'
#
#
# def test_notifier_send_oq_up_throughput_delay():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['throughput_delay'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysUpThroughputDelay'
#
#
# def test_notifier_send_oq_down_throughput_loss():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['throughput_loss'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysDownLossThroughput'
#
#
# def test_notifier_send_oq_up_throughput_loss():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['throughput_loss'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrend15daysUpLossThroughput'
#
#
# def test_notifier_send_oq_up_all_reasons_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['all'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateUpAll'
#
#
# def test_notifier_send_oq_down_all_reasons_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['all'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateDownAll'
#
#
# def test_notifier_send_oq_up_loss_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['loss'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateUpLoss'
#
#
# def test_notifier_send_oq_down_loss_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['loss'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateDownLoss'
#
#
# def test_notifier_send_oq_up_delay_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['delay'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateUpDelay'
#
#
# def test_notifier_send_oq_down_delay_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['delay'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateDownDelay'
#
#
# def test_notifier_send_oq_up_throughput_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['throughput'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateUpThroughput'
#
#
# def test_notifier_send_oq_down_throughput_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['throughput'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateDownThroughput'
#
#
# def test_notifier_send_oq_down_loss_delay_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['loss_delay'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateDownLossDelay'
#
#
# def test_notifier_send_oq_up_loss_delay_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['loss_delay'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateUpLossDelay'
#
#
# def test_notifier_send_oq_down_throughput_delay_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['throughput_delay'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message ==
# 'MidtermOQCNTrendImmediateDownThroughputDelay'
#
#
# def test_notifier_send_oq_up_throughput_delay_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['throughput_delay'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateUpThroughputDelay'
#
#
# def test_notifier_send_oq_down_throughput_loss_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['throughput_loss'],
#         send=True,
#         direction=DOWN)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateDownLossThroughput'
#
#
# def test_notifier_send_oq_up_throughput_loss_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE,
#             TEST_DAYS_TO_REPORT
#         ],
#         trend_type='immediate',
#         reasons=reasons['throughput_loss'],
#         send=True,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message == 'MidtermOQCNTrendImmediateUpLossThroughput'
#
#
# def test_notifier_no_send_lack_data():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['loss'],
#         send=False,
#         direction=UP)
#     notifier.input(mock_data)
#     assert notifier._message is None
#
#
# def test_notifier_no_send_lack_data_immediate():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=None,
#         trend_type=None,
#         reasons=None,
#         send=False,
#         direction=None)
#     notifier.input(mock_data)
#     assert notifier._message is None
#
#
# def test_notifier_no_send_nooq_data():
#     notifier = Notifier(aid_service_connection=mock_aid_service)
#     mock_data = generate_mock(
#         message_details=[
#             TEST_DT, TEST_TS, TEST_SCORES, TEST_APPID, TEST_CHANGE
#         ],
#         trend_type='monthly',
#         reasons=reasons['loss'],
#         send=False,
#         direction=UP)
#     # NOTE: Now the dict with data is sent, but OQ is not there.
#     #       Message should not be sent in this scenario.
#     del mock_data['shortterm']['complex_dm_output'][
#         'objective_quality_v35_average']  # noqa
#     del mock_data['midterm']['complex_dm_output'][
#         'objective_quality_v35_average']  # noqa
#     notifier.input(mock_data)
#     assert notifier._message is None
