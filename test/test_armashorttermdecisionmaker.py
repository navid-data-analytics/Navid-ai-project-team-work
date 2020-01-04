"""Test predictor class."""
from src.decisionmakers import ShortArmaDecisionMaker as ArmaDecisionMaker #noqa
from .MockMessageClient import MockMessageClient
import datetime
import numpy as np
import pandas as pd

forecast = np.asarray([
    19.98606423, 19.79212946, 19.53536529, 19.23137634, 18.90776768,
    18.57794709, 18.25184225, 17.93533362, 17.63204103, 17.3439433,
    17.07196531, 16.81633183, 16.57681997, 16.35292655, 16.14398162,
    15.94922416, 15.76785233, 15.5990563,  15.44203932, 15.29603079,
    15.16029403, 15.03413043, 14.91688124, 14.80792783, 14.70669093,
    14.61262929, 14.52523794, 14.44404634, 14.36861634, 14.29854024,
])

reversed_forecast = forecast[::-1]

rate_of_change = -28.507298803899193

time_stub = pd.Timestamp(datetime.datetime.utcnow().replace(year=2007, month=7,
                                                            day=2, hour=0,
                                                            minute=0,
                                                            second=0))

arma_output_growth = (reversed_forecast, -rate_of_change, time_stub)

arma_output_decline = (forecast, rate_of_change, time_stub)


def initialize_adm(app_id, decisions):
    mock_aid_service_connection = MockMessageClient()
    adm = ArmaDecisionMaker(app_id, metric='conferences_terminated',
                            aid_service_connection=mock_aid_service_connection)
    adm.output.connect(decisions.append)
    return adm


def trigger_adm(adm, trigger):
    adm.input(trigger)


def test_adm_pushes_datetime(decisions=list(),
                             triggers=arma_output_growth):
    adm = initialize_adm(app_id=1, decisions=decisions)
    trigger_adm(adm, triggers)
    assert isinstance(decisions[0]['time_model'], datetime.date)


def test_adm_pushes_verdict(decisions=list(),
                            triggers=arma_output_growth):
    adm = initialize_adm(app_id=1, decisions=decisions)
    trigger_adm(adm, triggers)
    assert isinstance(decisions[0]['decision'], int)


def test_adm_pushes_change_rate(decisions=list(),
                                triggers=arma_output_decline):
    adm = initialize_adm(app_id=2, decisions=decisions)
    trigger_adm(adm, triggers)
    assert isinstance(decisions[0]['rate_of_change'], float)


def test_adm_decides_growth(decisions=list(),
                            triggers=arma_output_growth,
                            growth=1):
    adm = initialize_adm(app_id=3, decisions=decisions)
    trigger_adm(adm, triggers)
    assert decisions[0]['decision'] == growth


def test_adm_decides_decline(decisions=list(),
                             triggers=arma_output_decline,
                             decline=-1):
    adm = initialize_adm(app_id=4, decisions=decisions)
    trigger_adm(adm, triggers)
    assert decisions[0]['decision'] == decline


def test_adm_messages_growth(decisions=list(),
                             triggers=arma_output_growth):
    adm = initialize_adm(app_id=5, decisions=decisions)
    trigger_adm(adm, triggers)
    assert 'grow' in adm.message


def test_adm_messages_decline(decisions=list(),
                              triggers=arma_output_decline):
    adm = initialize_adm(app_id=6, decisions=decisions)
    trigger_adm(adm, triggers)
    assert 'decline' in adm.message


def test_adm_receives_none(decisions=list(),
                           triggers=arma_output_decline):
    adm = initialize_adm(app_id=7, decisions=decisions)
    trigger_adm(adm, [None, None, None])
    assert 'Not enough data' in adm.message
