"""Predict an anomaly score on the data that was provided."""
import numpy as np
from src.models import Model
from src.utils import measure_time_metric
from collections import deque
import logging
from prometheus_client import Gauge
from scipy.stats import ttest_ind

model_output_prometheus = Gauge('Model_output', 'appIDs',
                                ['app_id', 'metric_name'])
model_average_traffic_prometheus = Gauge('model_average_traffic', 'AppIDs',
                                         ['app_id', 'metric_name'])
model_old_average_traffic_prometheus = Gauge(
    'model_old_average_traffic', 'AppIDs', ['app_id', 'metric_name'])

model_t_value_prometheus = Gauge('model_t_value', 'AppIDs',
                                 ['app_id', 'metric_name'])
model_p_value_prometheus = Gauge('model_p_value', 'AppIDs',
                                 ['app_id', 'metric_name'])

app_traffic_prometheus = Gauge('app_traffic', 'AppIDs',
                               ['app_id', 'metric_name'])

logger = logging.getLogger('root')


class TrendDetector(Model):
    """
    A Model object is responsible for determining whether the data fed to it
    by preprocessor object contains anomaly or not.
    The workflow should look like this:
    - First model is initiated by configuration parameters received as dict.
    - Model object receives a DataFrame as input and first separates the time
     and value  columns into two variables.
    - The output of the model is a list of detected_points and the dateTime.
    - 'detected_points' is a list of zeros and ones where the ones indicate
    the positions that change was detected. e.g., if detected point is like:
    [(Timestamp('2016-02-01 00:00:00'), 0),
    (Timestamp('2016-02-02 00:00:00'), 0),
    (Timestamp('2016-02-03 00:00:00'), 0),
    (Timestamp('2016-02-04 00:00:00'), 1)],
     it means the on 2016-02-04 something abnormal has happened.
    """

    def __init__(self,
                 model_configuration,
                 app_id=None,
                 metric=None,
                 process=None):
        """Construct a new Model instance."""
        process = self._run if process is None else process
        super(TrendDetector, self).__init__(
            app_id=app_id, process=self._run, metric=metric)
        logger.debug(str(self) + ' Creating Model object')
        logger.debug(str(self) + ' Setting initial configuration of the model')
        self.data_index = 1
        self.variation = 0
        self.sliding_window_size = model_configuration['sliding_window_size']

        self.sliding_window = deque(maxlen=self.sliding_window_size)
        self.sliding_window_large = deque(maxlen=self.sliding_window_size * 2)
        self.t_value = 0
        self.p_value = 0
        self._no_traffic_days = 0
        self.population_flag = False
        self.overlappingLowIndex = \
            np.round(self.sliding_window_size/2) + self.sliding_window_size
        self.sliding_window_large_size =\
            self.sliding_window_size * 2
        self.overlappingUpIndex =\
            self.sliding_window_large_size +\
            np.round(self.sliding_window_size/2)-1
        self.confidence_threshold = model_configuration['confidence_threshold']
        logger.debug(str(self) + ' Model created')
        self._lower_index = int(self.sliding_window_size / 2)
        self._upper_index = int(self._lower_index + self.sliding_window_size)

    def _no_traffic_flag(self, new_signal_value):
        """
        Check if there have been more than two days without traffic.

        Arguments:
        - new_signal_value: int, number of traffic per day

        Returns:
        - int, 0 if traffic, 1 if no traffic for more than 3 days
        """
        logger.debug(repr(self) + 'Checking number of no-traffic days')
        if new_signal_value == 0:
            logger.debug(repr(self) + 'New value is a no-traffic day')
            self._no_traffic_days = self._no_traffic_days + 1
        else:
            logger.debug(repr(self) + 'Reset no-traffic day counter')
            self._no_traffic_days = 0

        if self._no_traffic_days >= 3:
            logger.debug(repr(self) + 'No traffic for 3 days or more')
            return 1
        else:
            logger.debug(repr(self) + 'Traffic detected')
            return 0

    def _populating_windows(self, new_signal_value):
        """
        Populate the sliding windows.

         If the windows need to have some overlapping,
          this function takes care of that

        Arguments:
        - new_signal_value: int, the value of time series
        """
        if self.data_index < self.sliding_window_large_size:
            logger.debug(
                repr(self) +
                'Appending new signal value to large sliding window')
            self.sliding_window_large.append(new_signal_value)
        if self.overlappingLowIndex <= self.data_index <= \
                self.overlappingUpIndex:
            logger.debug(
                repr(self) + 'Appending new signal value sliding window')
            self.sliding_window.append(new_signal_value)
            if self.data_index == self.overlappingUpIndex:
                logger.debug(repr(self) + 'Setting population_flag to True')
                self.population_flag = True
        self.data_index += 1

    def _execute(self, new_signal_value):
        """
        The statistical inference is implemented here.

        A two-tail t-distribution test
        (https://en.wikipedia.org/wiki/Welch%27s_t-test)
        has been applied from scipy library. The method
        looks for the significance of deviation between
        means of two distributions (two sliding windows)

        Arguments:
        - new_signal_value: int: value in current time

        Returns:
        - int, 1,0 or -1 for upward, no trend downward
        """
        self.sliding_window_large.append(self.sliding_window[int(
            self.sliding_window_size / 2)])
        self.sliding_window.append(new_signal_value)
        try:
            logger.debug(repr(self) + 'Executing T-Test')
            self.t_value, self.p_value = ttest_ind(
                self.sliding_window,
                self.sliding_window_large,
                equal_var=False)
            confidence_level = 1 - self.p_value
        except ValueError:
            logger.debug(
                repr(self) + 'Something went wrong, confidence level set to 0')
            self.t_value = "fail"
            self.p_value = "fail"
            confidence_level = 0
        self.data_index += 1
        result = self._predict(confidence_level)
        logger.debug(
            repr(self) + 'Notifying about the decision: {}'.format(result))
        return result

    def _predict(self, confidence_level):
        """
        Perform predition and return verdict.

        Arguments:
        - confidence_level: float

        Returns:
        - int: 1, 0 or -1 for upward, no trend downward
        """
        logger.debug(
            repr(self) + 'Confidence level: {}'.format(confidence_level))
        if (np.abs(confidence_level) > self.confidence_threshold) &\
                (self.t_value > 0):
            logger.debug(repr(self) + 'Detected change upward, sending 1')
            return 1
        elif (np.abs(confidence_level) > self.confidence_threshold) &\
                (self.t_value < 0):
            logger.debug(repr(self) + 'Detected change downward, sending -1')
            return -1
        else:
            logger.debug(repr(self) + 'No change detected, sending 0')
            return 0

    @measure_time_metric
    def _run(self, input_signal):
        """
        It is triggered upon arrival of each new data.

        If the two sliding windows have not been
        fully populated yet, it calls _populating_windows.
        Otherwise, it calls _execute.
        Arguments:
        - input_signal: each new input has a time and value
                        as a dataFrame object
        Returns:
        - time_model: Pandas datetime
        - value_model: integer, the decision of model, 1,0 or -1 for upward,
                       no trend or downward
        - app_id_model: int, the appID
        - no_traffic_flag: bool, if there have been 3 days with no traffic
        """
        logger.debug(repr(self) + '_run started')
        new_signal_value = input_signal.reset_index(drop=True).value[0]
        logger.debug(
            repr(self) + 'Value of latest signal: {}'.format(new_signal_value))
        current_time = input_signal.reset_index(drop=True).time[0]
        logger.debug(repr(self) + 'Current time: {}'.format(current_time))
        if self._no_traffic_flag(new_signal_value) == 1:
            logger.debug(repr(self) + 'No recent traffic detected')
            no_traffic = True
        else:
            logger.debug(repr(self) + 'Recent traffic detected')
            no_traffic = False

        if self.population_flag:
            logger.debug(repr(self) + 'Population flag set to True, proceed')
            model_decision = self._execute(new_signal_value)
            avg_traffic = np.mean(self.sliding_window)
            logger.debug(
                repr(self) +
                'Average traffic calculated: {}'.format(avg_traffic))
            old_average_traffic = np.mean(
                list(self.sliding_window_large)[self._lower_index:self.
                                                _upper_index])
        else:
            logger.debug(repr(self) + 'Population flag set to False')
            self._populating_windows(new_signal_value)
            model_decision = 0
            avg_traffic = 0
            old_average_traffic = 0
            logger.debug(
                repr(self) + 'Model decision and average traffic set to 0')

        self._postprocess(model_decision, avg_traffic, old_average_traffic,
                          new_signal_value)

        model_output = {
            'time_model': current_time,
            'value_model': model_decision,
            'app_id_model': self._app_id,
            'avg_traffic': avg_traffic,
            'old_average_traffic': old_average_traffic,
            'no_traffic_flag': no_traffic
        }
        logger.debug(repr(self) + 'Output of Model'.format(model_output))
        return model_output

    def _postprocess(self, model_decision, avg_traffic, old_average_traffic,
                     new_signal_value):
        """Push decision and traffic info to prometheus."""
        logger.debug(
            repr(self) +
            'Send model_decision and average traffic to prometheus')
        model_output_prometheus.labels(self._app_id,
                                       self.metric).set(model_decision)
        model_average_traffic_prometheus.labels(self._app_id,
                                                self.metric).set(avg_traffic)
        app_traffic_prometheus.labels(self._app_id,
                                      self.metric).set(new_signal_value)
        model_old_average_traffic_prometheus.labels(
            self._app_id, self.metric).set(old_average_traffic)

        # Clean up 'false' instances
        p_value, t_value = [
            val if not isinstance(val, str) else 0
            for val in [self.p_value, self.t_value]
        ]
        model_p_value_prometheus.labels(self._app_id, self.metric).set(p_value)
        model_t_value_prometheus.labels(self._app_id, self.metric).set(t_value)

        logger.debug(
            repr(self) + 'Decision and average traffic sent to prometheus')

    def __repr__(self):
        """Use repr(self) for DEBUG, ERROR, CRITICAL logging levels."""
        return '{} - {} - Sliding_window {}: '.format(  #noqa
            self.app_id, self.metric, self.sliding_window_size)

    def __str__(self):
        """Use str(self) for INFO, WARNING logging levels."""
        return '{}: '.format(self.app_id)
