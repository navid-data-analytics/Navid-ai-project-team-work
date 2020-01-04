import time
import logging
from prometheus_client import Gauge

logger = logging.getLogger('root')
prometheus_time_gauge = Gauge('execution_time_of_last_trigger',
                              'time_elapsed to execute', ['function_name'])
prometheus_metric_time_gauge = Gauge('execution_time_of_last_trigger_metric',
                                     'time_elapsed to execute',
                                     ['function_name', 'metric_name',
                                      'app_id'])


def measure_time(function):
    def _measure(*args, **kwargs):
        start = time.perf_counter()
        result = function(*args, **kwargs)
        end = time.perf_counter()
        execution_time = end - start
        logger.debug('{}.{} executed in {} seconds'.format(
            function.__module__, function.__name__, execution_time))
        prometheus_time_gauge.labels(function.__qualname__).set(execution_time)
        return result
    return _measure


def measure_time_metric(function):
    def _measure(*args, **kwargs):
        start = time.perf_counter()
        result = function(*args, **kwargs)
        end = time.perf_counter()
        execution_time = end - start
        logger.debug('{}.{} for {} {} executed in {} seconds'.format(
            function.__module__, function.__name__, args[0].app_id,
            args[0].metric, execution_time))
        prometheus_metric_time_gauge.labels(function.__qualname__,
                                            args[0].metric,
                                            args[0].app_id).set(execution_time)
        return result
    return _measure
