"""Tests the execution of .(all components connected)."""
import datetime
import logging
from src.components import DatePackage, Broadcaster, Joiner
from src.pipeline import DateQueue, DataFetcher, Notifier, MetricPipeline
from src.pipeline import AppPipelineVolume, AppPipelineOQ, AppPipelineRTT
from src.pipeline import AppPipelineThroughputEffect, AppPipelineDelayEffect
from src.pipeline import AppPipelineLossEffect, AppPipelineRTTShortTerm
from src.pipeline import AppPipelineVolumeShortTerm, AppPipelineOQShortTerm
from src.decisionmakers import ComplexOQDecisionMaker
from src.Grpc.CrsClient import CrsClient
from itertools import product
import src.utils.constants as constants
import time
import numpy as np

np.random.seed(0)
logger = logging.getLogger('root')

app_ids = [427974000]

supported_metrics = {
    # SHORTTERM
    'shortterm': {
        'conferences_terminated': AppPipelineVolumeShortTerm,
        'objective_quality_v35_average': AppPipelineOQShortTerm,
        'rtt_average': AppPipelineRTTShortTerm,
        'delayEffect': AppPipelineDelayEffect,
        'throughputEffect': AppPipelineThroughputEffect,
        'lossEffect': AppPipelineLossEffect,
    },
    # MIDTERM
    'midterm': {
        'conferences_terminated': AppPipelineVolume,
        'objective_quality_v35_average': AppPipelineOQ,
        'rtt_average': AppPipelineRTT,
        'delayEffect': AppPipelineDelayEffect,
        'throughputEffect': AppPipelineThroughputEffect,
        'lossEffect': AppPipelineLossEffect,
    }
}

supported_terms = ('midterm', 'shortterm')

metrics = [
    'delayEffect', 'throughputEffect', 'lossEffect', 'conferences_terminated',
    'objective_quality_v35_average', 'rtt_average'
]

mock_configuration = {
    'detection': {
        'sliding_window_size': 30,
        'confidence_threshold': 0.8
    },
    'prediction': {
        'history_size': 24,
        'seasonality_period': 7,
        'arma_order': {
            427974000: (3, 0)
        },
    },
    'fluctuation': {
        427974000: {
            'thresholds': (100, 200)
        }
    },
}

config = {
    'models': {
        # SHORTTERM
        'shortterm': {
            'conferences_terminated': mock_configuration,
            'objective_quality_v35_average': mock_configuration,
            'rtt_average': mock_configuration,
            'delayEffect': mock_configuration,
            'throughputEffect': mock_configuration,
            'lossEffect': mock_configuration,
        },
        # MIDTERM
        'midterm': {
            'conferences_terminated': mock_configuration,
            'objective_quality_v35_average': mock_configuration,
            'rtt_average': mock_configuration,
            'delayEffect': mock_configuration,
            'throughputEffect': mock_configuration,
            'lossEffect': mock_configuration,
        }
    }
}

mock_data = {
    'conferences_terminated': 5,
    'objective_quality_v35_average': 2.0,
    'rtt_average': 150,
    'delayEffect': 1,
    'throughputEffect': 1,
    'lossEffect': 2,
}


class mock_service:
    def __init__(self):
        pass

    def Aggregate(self):
        pass


def mock_getService(stub):
    return mock_service()


def mock_send(method, request, name, reliable=False):
    return mock_data, None


def mock_aggregate(appID, from_dt, to_dt):
    return mock_data


def create_components(metrics, app_ids, model_conf=config):
    """
    Output collector to check how much data is pushed through,
    Notifier to see if the pipeline does not crash.
                                   OC
    DQ---DF---MetricPipeline--BC-<
                                   Nofifier
    """

    datequeue = DateQueue()

    client = CrsClient('notexistent:5432')
    client._connection_timeout = 0
    client._max_retries = 1
    client.getService = mock_getService
    client.send = mock_send
    client.Aggregate = mock_aggregate
    datafetcher = DataFetcher(metrics, app_ids, client)
    metricpipeline = {
        '{}_{}'.format(metric, term): MetricPipeline(
            supported_metrics[term][metric], app_ids,
            config['models'][term][metric])
        for metric, term in product(metrics, supported_terms)
    }

    metricjoiner = Joiner(
            process=lambda values: {'{}_{}'.format(
                metric, term): values['{}_{}'.format(metric, term)]
                for metric, term in product(metrics,
                                            supported_terms)})

    complex_decisionmaker_oq = ComplexOQDecisionMaker(app_ids=app_ids)
    fork = Broadcaster()
    notifier = Notifier()
    output_collector = list()
    return (datequeue, datafetcher, metricpipeline, metricjoiner,
            complex_decisionmaker_oq, fork, notifier, output_collector)


def initialize_components(components, metrics, app_ids):
    datequeue, datafetcher, metricpipeline, metricjoiner, \
        complex_decisionmaker_oq, fork, notifier, \
        output_collector = components
    datequeue.output.connect(datafetcher.input)
    for metric in metrics:
        for term in supported_terms:
            datafetcher.get_output('{}_{}'.format(metric, term)).connect(
                metricpipeline['{}_{}'.format(metric, term)].input)
            metricpipeline['{}_{}'.format(metric, term)].output.connect(
                metricjoiner.get_input('{}_{}'.format(metric, term)))
    metricjoiner.output.connect(fork.input)
    fork.output.connect(notifier.input)
    # complex_decisionmaker_oq.output.connect(notifier.input)
    fork.output.connect(output_collector.append)


def calculate_time_details(start, end, number_of_triggers):
    start_timestamp = time.mktime(start.timetuple()) * constants.CONVERT_KILO
    end_timestamp = time.mktime(end.timetuple()) * constants.CONVERT_KILO
    increment = (end_timestamp - start_timestamp) / number_of_triggers
    assert start_timestamp + increment * number_of_triggers == end_timestamp
    return start_timestamp, increment


def create_date_packages(start, increment, number_of_triggers):
    return [start + increment * _ for _ in range(number_of_triggers + 1)]


def trigger_datequeue(datequeue, trigger_times):
    packages = [
        DatePackage(start=trigger_times[i], end=trigger_times[i + 1])
        for i in range(len(trigger_times) - 1)
    ]
    [datequeue.input(p) for p in packages]


def mock_trigger(datequeue, start, end, number_of_triggers=1):
    start_timestamp, increment = calculate_time_details(
        start, end, number_of_triggers)
    trigger_times = create_date_packages(start_timestamp, increment,
                                         number_of_triggers)
    trigger_datequeue(datequeue, trigger_times)


def check_result(output_collector, number_of_triggers):
    result = len(output_collector)
    logger.debug(output_collector)
    assert result == number_of_triggers, 'Expected: {}, Actual: {}'.format(
        number_of_triggers, result)


def _test_pipeline(number_of_data_intervals,
                   start_time=datetime.datetime(2018, 2, 10, 0, 0, 0, 0,
                                                datetime.timezone.utc),
                   end_time=datetime.datetime(2018, 2, 10, 1, 0, 0, 0,
                                              datetime.timezone.utc)):

    components = create_components(metrics, app_ids, model_conf=config)
    initialize_components(components, metrics, app_ids)

    datequeue = components[0]
    output_collector = components[-1]
    mock_trigger(datequeue, start_time, end_time, number_of_data_intervals)
    datequeue.shutdown()
    datequeue.thread.join()
    check_result(output_collector, number_of_data_intervals)


def test_detection_one_interval(number_of_data_intervals=1):
    _test_pipeline(number_of_data_intervals)


def test_detection_two_intervals(number_of_data_intervals=2):
    _test_pipeline(number_of_data_intervals)


def test_detection_ten_intervals(number_of_data_intervals=10):
    _test_pipeline(number_of_data_intervals)


def test_detection_thirty_intervals(number_of_data_intervals=30):
    # NOTE: Here the arma history is filled and uses decomposed signal,
    #       with seasonality = 7 it needs to cut 7 days from the output
    expected_outputs = number_of_data_intervals - constants.WEEKLY
    _test_pipeline(expected_outputs)
