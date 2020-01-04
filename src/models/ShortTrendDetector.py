"""Predict an anomaly score on the data that was provided."""
import numpy as np
from src.models import TrendDetector
from src.utils import measure_time_metric
from collections import deque
import logging
from prometheus_client import Gauge
import scipy.stats as stats

model_output_prometheus = Gauge('shortterm_model_output', 'appIDs',
                                ['app_id', 'metric_name'])
model_average_traffic_prometheus = Gauge('shortterm_model_average_traffic',
                                         'AppIDs', ['app_id', 'metric_name'])
model_old_average_traffic_prometheus = Gauge(
    'shortterm_old_average_traffic', 'AppIDs', ['app_id', 'metric_name'])
app_traffic_prometheus = Gauge('shortterm_app_traffic', 'AppIDs',
                               ['app_id', 'metric_name'])

logger = logging.getLogger('root')


class ShortTrendDetector(TrendDetector):
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

    def __init__(self, model_configuration, app_id=None, metric=None):
        """Construct a new Model instance."""
        super(ShortTrendDetector, self).__init__(
            model_configuration=model_configuration,
            app_id=app_id,
            process=self._run,
            metric=metric)
        self.half_size = model_configuration['sliding_window_size']
        self.sliding_window = deque(maxlen=self.half_size * 2)
        logger.debug(str(self) + ' Model created')

    def _preprocess(self, input_signal):
        new_signal_value = input_signal.reset_index(drop=True).value[0]
        logger.debug(
            repr(self) + 'Value of latest signal: {}'.format(new_signal_value))
        current_time = input_signal.reset_index(drop=True).time[0]
        logger.debug(repr(self) + 'Current time: {}'.format(current_time))
        no_traffic = self._no_traffic_flag(new_signal_value) == 1
        return new_signal_value, current_time, no_traffic

    def _populate_windows(self, new_signal_value):
        """
        Populate the sliding windows.

         If the windows need to have some overlapping,
          this function takes care of that

        Arguments:
        - new_signal_value: int, the value of time series

        Returns:
        - boolean, populated flag
        """
        window_populated = len(self.sliding_window) >= self.half_size * 2
        logger.debug("Window populated {}".format(window_populated))
        logger.debug("sliding window {}, {}".format(
            len(self.sliding_window), self.half_size * 2))
        self.sliding_window.append(new_signal_value)
        return window_populated

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
        - avg_traffic - current average traffic
        - old_avg_traffic - previous average traffic
        """
        current_wnd = list(self.sliding_window)[self.half_size:]
        previous_wnd = list(self.sliding_window)[:self.half_size]
        avg_traffic = np.mean(current_wnd)
        logger.debug(
            repr(self) + 'Average traffic calculated: {}'.format(avg_traffic))
        logger.debug("Sliding window {}".format(self.sliding_window))
        logger.debug("Current window {}".format(current_wnd))
        logger.debug("Previous window {}".format(previous_wnd))
        old_avg_traffic = np.mean(previous_wnd)

        logger.debug(repr(self) + 'Executing T-Test')
        self.t_value, self.p_value = stats.ttest_rel(current_wnd, previous_wnd)
        confidence_level = 1 - self.p_value

        result = self._predict(confidence_level)
        logger.debug(
            repr(self) + 'Notifying about the decision: {}'.format(result))
        return result, avg_traffic, old_avg_traffic

    @measure_time_metric
    def _run(self, input_signal):
        """
        It is triggered upon arrival of each new data.

        If the two sliding windows have not been
        fully populated yet, it calls _populate_windows.
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
        decision, avg_traffic, old_avg_traffic = 0, 0, 0

        new_signal_value, time, no_traffic = self._preprocess(input_signal)

        window_populated = self._populate_windows(new_signal_value)
        logger.debug('Window populated? {}'.format(window_populated))
        if window_populated:
            logger.debug("Run _execute")
            decision, avg_traffic, old_avg_traffic = self._execute(
                new_signal_value)

        self._postprocess(decision, avg_traffic, old_avg_traffic,
                          new_signal_value)

        model_output = {
            'time_model': time,
            'value_model': decision,
            'app_id_model': self._app_id,
            'avg_traffic': avg_traffic,
            'old_average_traffic': old_avg_traffic,
            'no_traffic_flag': no_traffic
        }
        logger.debug(repr(self) + 'Output of Model'.format(model_output))
        return model_output
