"""Test trend decision maker class. The input to model is
a signal where the decisionmaker should detect
 one monthly upward trend """
from src.decisionmakers import DecisionMakerMidTerm
from .MockMessageClient import MockMessageClient
import logging
import pandas as pd
import datetime

logger = logging.getLogger('root')
logger.setLevel('DEBUG')
SUPPORTED_METRICS = ['conferences_terminated']


def create_inputs(vals_list, start_date):
    """
    This method creates some inputs for testing
     the performance of Decision maker
    :param vals_list: a list of integers. [output of model,
     old average traffic, current average traffic]
    :param start_date: datetime object
    :return: list of dicts. Each dict is equivalent with
     the input of DM coming form model
    """
    inputs = []
    for i in range(len(vals_list)):
        val, old_avg_traffic, avg_traffic = vals_list[i]
        input = {'time_model': pd.to_datetime(
            start_date + datetime.timedelta(days=i)
        ),
            'app_id_model': 1234,
            'no_traffic_flag': False,
            'avg_traffic': avg_traffic,
            'old_average_traffic': old_avg_traffic,
            'value_model': val}
        inputs.append(input)
    return inputs


def create_dm(metric_name):
    mock_aid_service_connection = MockMessageClient()
    dm_out = []
    trenddecisionmaker = DecisionMakerMidTerm(
                          metric=metric_name,
                          app_id=123,
                          aid_service_connection=mock_aid_service_connection)
    trenddecisionmaker.output.connect(dm_out.append)

    return trenddecisionmaker, dm_out


def check_for_immediate_decisions(inputs=[],
                                  decisions=[],
                                  start_date=datetime.datetime(
                           2018, 7, 10, 0, 0, 0,
                           tzinfo=datetime.timezone.utc)):
    """
    This method is used for testing the immediate
     decision of DM
    :param inputs: list. [val= int and output of model,
     old_avg_traffic = (integer) old average traffic,
      avg_traffic= (integer) current average traffic]
    :param decisions: list of int. The expected input of model
    :param start_date: datetime object.
    :return: bool. True if expected output of model is met,
     False otherwise
    """
    inputs = create_inputs(inputs, start_date)
    for metric_name in SUPPORTED_METRICS:
        dm, dm_out = create_dm(metric_name)
        for point in inputs:
            dm.input(point)

        for i in range(len(decisions)):
            assert decisions[i] == dm_out[i]['trend_immediate_decision']
            if inputs[i]['time_model'].day != 16:
                assert dm_out[i]['trend_monthly_decision'] != 1


def check_for_monthly_decisions(inputs=[],
                                decisions=[],
                                start_date=datetime.datetime(
                            2018, 7, 10, 0, 0, 0,
                            tzinfo=datetime.timezone.utc)):
    """
    This method is used for testing the monthly decision of DM
    :param inputs: list, [val= int and output of model,
     old_avg_traffic = (integer) old average traffic,
     avg_traffic= (integer) current average traffic]
    :param decisions: list of integers. The expected output
     of model
    :param start_date: datetime object.
    :return: bool. True if expected output of model is met,
     False otherwise
    """
    inputs = create_inputs(inputs, start_date)
    for metric_name in SUPPORTED_METRICS:
        dm, dm_out = create_dm(metric_name)
        for point in inputs:
            logger.info(point)
            dm.input(point)

        for i in range(len(decisions)):
            assert decisions[i] == dm_out[i]['trend_monthly_decision']
            if inputs[i]['time_model'].day == 16:
                assert 1 != dm_out[i]['trend_immediate_decision']


def test_reports_immediate_trend_up(caplog):
    """
    This function checks both the output of DM which is 0,1 or -1
     using check_for_immediate_decisions() method
     and also if message was sent to AID-E. The test here is:
     there are two trendy days detected by model on 10 and 11 of July.
      It is expected that DM sends message only on the first day
    :param check_for_immediate_decisions: A method which checks if
     the return value of DM makes sense
    :param caplog: From pytest-catchlog Library.
     Used for log capturing
    :return:bool, True if message was sent, False otherwise
    """
    check_for_immediate_decisions(
        inputs=[
            [0, 100, 200],
            [1, 100, 200],
            [1, 101, 205]
        ],
        decisions=[
            0, 1, 0
        ],
        start_date=datetime.datetime(2018, 7, 10, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )
    assert "immediate upwards trend was sent." in caplog.text


def test_reports_immediate_trend_down(caplog):
    """
    The test here is:
     there is only one downward trend on 11th of July,
      it is expected that DM sends message message only on 11th
    :param check_for_immediate_decisions: A method which checks if
     the return value of DM makes sense
    :param caplog: From pytest-catchlog Library.
     Used for log capturing
    :return:bool, True if message was sent, False otherwise
    """
    check_for_immediate_decisions(
        inputs=[
            [0, 100, 200],
            [-1, 101, 50],
            [-1, 101, 50]
        ],
        decisions=[
            0, -1, 0
        ],
        start_date=datetime.datetime(2018, 7, 10, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )
    assert "immediate downwards trend was sent." in caplog.text


def test_reports_monthly_trend_up(caplog):
    """
    The test here is:
     there is two upward trendy days on 13th and 14th of July,and
      there is no trend on 15th and 16th while 16th is
       the monthly report day. It is expected that an upward trend
        monthly report is sent on 16th
    :param check_for_immediate_decisions: A method which checks if
     the return value of DM makes sense
    :param caplog: From pytest-catchlog Library.
     Used for log capturing
    :return:bool, True if message was sent, False otherwise
    """
    check_for_monthly_decisions(
        inputs=[
            [1, 100, 200],
            [1, 101, 205],
            [0, 101, 205],
            [0, 101, 205],
        ],
        decisions=[
            0, 0, 0, 1
        ],
        start_date=datetime.datetime(2018, 7, 13, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )
    assert "monthly upwards trend was sent." in caplog.text


def test_reports_monthly_trend_down(caplog):
    """
    The test here is:
     there is two downward trendy days on 14th and 15th of July, and
      there is no trend on 16th while 16th is
       the monthly report day. It is expected that a downward trend
        monthly report is sent on 16th
    :param check_for_immediate_decisions: A method which checks if
     the return value of DM makes sense
    :param caplog: From pytest-catchlog Library.
     Used for log capturing
    :return:bool, True if message was sent, False otherwise
    """
    check_for_monthly_decisions(
        inputs=[
            [0, 100, 20],
            [-1, 101, 50],
            [-1, 102, 55],
            [0, 101, 56]
        ],
        decisions=[
            0, 0, 0, -1
        ],
        start_date=datetime.datetime(2018, 7, 13, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )
    assert "monthly downwards trend was sent." in caplog.text


def test_for_no_monthly_trend():
    """
    The test here is:
     there has benn no trend and monthly report function should not
      report anything
    :param check_for_immediate_decisions: A method which checks if
     the return value of DM makes sense
    """
    check_for_monthly_decisions(
        inputs=[
            [0, 100, 100],
            [0, 101, 100],
            [0, 102, 100],
            [0, 101, 101]
        ],
        decisions=[
            0, 0, 0, 0
        ],
        start_date=datetime.datetime(2018, 7, 13, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )


def test_for_no_immediate_trend():
    """
    The test here is:
     there has benn no trend, and immediate decision maker should not
      report anything
    :param check_for_immediate_decisions: A method which checks if
     the return value of DM makes sense
    """
    check_for_immediate_decisions(
        inputs=[
            [0, 100, 100],
            [0, 101, 100],
            [0, 102, 100],
            [0, 101, 101]
        ],
        decisions=[
            0, 0, 0, 0
        ],
        start_date=datetime.datetime(2018, 7, 13, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )


def test_for_supressing_the_downward_message():
    """
    The test here is:
     even if there is one downward trend on 11th of July,
      it is expected that DM does not send any message
       because there has not been any decline ratio: (102-101)>0
    :param check_for_immediate_decisions: A method which checks if
     the return value of DM makes sense
    """
    check_for_immediate_decisions(
        inputs=[
            [0, 100, 200],
            [-1, 101, 102]
        ],
        decisions=[
            0, 0
        ],
        start_date=datetime.datetime(2018, 7, 10, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )


def test_for_supressing_the_upward_message():
    """
    The test here is:
     even if there is one upward trend on 11th of July,
      it is expected that DM does not send any message
       because there has not been any growth ratio: (100-99)>0
    :param check_for_immediate_decisions: A method which checks if
     the return value of DM makes sense
    """
    check_for_immediate_decisions(
        inputs=[
            [0, 100, 100],
            [1, 100, 99]
        ],
        decisions=[
            0, 0
        ],
        start_date=datetime.datetime(2018, 7, 10, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )


def test_down_to_upwards_immediate_trend(caplog):
    """
    The test here is:
     If the trend changes from downward to upward it should send
       both messages.
       There is always a "no trend" in between as the average changes slowly.

    :param check_for_immediate_decisions: A method which checks if
     the return value of DM makes sense
    """
    check_for_immediate_decisions(
        inputs=[
            [-1, 100, 80],
            [0, 100, 100],
            [1, 100, 120]
        ],
        decisions=[
            -1, 0, 1
        ],
        start_date=datetime.datetime(2018, 7, 10, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )
    assert "immediate downwards trend was sent." in caplog.text
    assert "immediate upwards trend was sent." in caplog.text


def test_up_to_downwards_immediate_trend(caplog):
    """
    The test here is:
     If the trend changes from upward to downward it should send
       both messages.
       There is always a "no trend" in between as the average changes slowly.

    :param check_for_immediate_decisions: A method which checks if
     the return value of DM makes sense
    """
    check_for_immediate_decisions(
        inputs=[
            [1, 100, 120],
            [0, 100, 100],
            [-1, 100, 80]
        ],
        decisions=[
            1, 0, -1
        ],
        start_date=datetime.datetime(2018, 7, 10, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )
    assert "immediate upwards trend was sent." in caplog.text
    assert "immediate downwards trend was sent." in caplog.text


def test_monthly_down_suppression(caplog):
    """
    The test here is:
     If there was a trend within the last 15 days, but on monthly
     report day we do not see a change in averages, don't report.
    :param check_for_monthly_decisions: A method which checks if
     the return value of DM makes sense
    """
    check_for_monthly_decisions(
        inputs=[
            [0, 100, 20],
            [-1, 101, 50],
            [-1, 102, 55],
            [0, 101, 101]
        ],
        decisions=[
            0, 0, 0, 0
        ],
        start_date=datetime.datetime(2018, 7, 13, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )
    assert "monthly downwards trend was sent." not in caplog.text


def test_monthly_up_suppression(caplog):
    """
    The test here is:
     If there was a trend within the last 15 days, but on monthly
     report day we do not see a change in averages, don't report.
    :param check_for_monthly_decisions: A method which checks if
     the return value of DM makes sense
    """
    check_for_monthly_decisions(
        inputs=[
            [0, 100, 20],
            [1, 101, 120],
            [1, 102, 122],
            [0, 101, 101]
        ],
        decisions=[
            0, 0, 0, 0
        ],
        start_date=datetime.datetime(2018, 7, 13, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )
    assert "monthly upwards trend was sent." not in caplog.text


def test_immediate_down_suppression(caplog):
    """
    The test here is:
     If there was an immediate trend but no right-direction
     change in the averages, don't report.
    """
    check_for_monthly_decisions(
        inputs=[
            [0, 100, 20],
            [-1, 101, 102],
        ],
        decisions=[
            0, 0
        ],
        start_date=datetime.datetime(2018, 7, 10, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )
    assert "immediate downwards trend was sent." not in caplog.text


def test_immediate_up_suppression(caplog):
    """
    The test here is:
     If there was an immediate trend but no right-direction
     change in the averages, don't report.
    """
    check_for_monthly_decisions(
        inputs=[
            [0, 100, 20],
            [1, 102, 101],
        ],
        decisions=[
            0, 0
        ],
        start_date=datetime.datetime(2018, 7, 10, 0, 0, 0,
                                     tzinfo=datetime.timezone.utc)
    )
    assert "immediate downwards trend was sent." not in caplog.text


# TODO both up and downward trend within last 15 days
# currently behavior of DM is not well defined,
# as this should not happen in midterm
