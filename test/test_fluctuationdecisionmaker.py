from src.decisionmakers import FluctuationDecisionMaker
from .MockMessageClient import MockMessageClient
import datetime
import pandas as pd

SUPPORTED_METRICS = ('conferences_terminated', 'objective_quality_v35_average',
                     'rtt_average')


def create_inputs(vals_list):
    return [{'time_model': pd.to_datetime(
                datetime.datetime(2018, 7, 10, 11, 22, 3,
                                  tzinfo=datetime.timezone.utc)),
             'app_id_model': 'dummy',
             'value_model': val} for val in vals_list]


def setup_model(metric_name):
    mock_aid_service_connection = MockMessageClient()
    model_out = []
    model = FluctuationDecisionMaker(
                        metric=metric_name,
                        app_id=123,
                        aid_service_connection=mock_aid_service_connection)
    model.output.connect(model_out.append)
    return model, model_out


def check_message_type(inputs=[0], out_ind=0, message=None):
    inputs = create_inputs(inputs)
    for metric_name in SUPPORTED_METRICS:
        model, model_out = setup_model(metric_name)
        for point in inputs:
            model.input(point)
        if message is None:
            assert model_out[out_ind]['message'] is None
        else:
            assert model_out[out_ind]['message'] == model.messages[message]


def test_if_reports_abnormal_increase():
    check_message_type([1], 0, 1)


def test_if_reports_abnormal_decrease():
    check_message_type([-1], 0, -1)


def test_if_reports_decrease_to_norm():
    check_message_type([1, 0], 1, 0)


def test_if_reports_increase_to_norm():
    check_message_type([-1, 0], 1, 0)


def test_if_reports_none_if_stable_norm():
    check_message_type([0, 0], 1, None)


def test_if_reports_none_if_stable_high():
    check_message_type([1, 1], 1, None)


def test_if_reports_none_if_stable_low():
    check_message_type([-1, -1], 1, None)
