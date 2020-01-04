import numpy as np
from src.models import Model
from src.utils import running_mean_fast, running_fn_norm, measure_time_metric
from collections import deque
import logging
from prometheus_client import Gauge
from math import isclose
import src.utils.constants as constants

model_output_prometheus = Gauge('Fluctuation_model_output', 'appIDs',
                                ['app_id', 'metric_name'])
model_metric_prometheus = Gauge('Fluctuation_metric_output', 'appIDs',
                                ['app_id', 'metric_name'])

logger = logging.getLogger('root')

PERIOD = 7  # weekly periodicity


class MidTermFluctuationDetector(Model):
    """
    Mid-term fluctuation detector component.

    Fluctuation detector determines if the mid-term data (~30 days)
    fluctuation level is normal, abnormally high, or abnormally low
    with the outputs of 0, 1, -1.
    """

    def __init__(self,
                 model_configuration,
                 app_id=1,
                 window_size=14,
                 metric=None):
        """Construct a MidTermFluctuationDetector object.

        Arguments:
        - model configuration: dict of dicts with appID:config
                               pairs, each config should have
                               'thresholds' key with a tuple of
                               (low_thresh, high_thresh)
        - app_id: int, self-explainatory
        """
        super(MidTermFluctuationDetector, self).__init__(
            app_id=app_id, process=self._run, metric=metric)
        self._window_size = window_size
        self._half_window = int(self._window_size / 2)
        self._memory_size = self._window_size * 3  # 3x windowing
        self.memory = deque(maxlen=self._memory_size)
        appconfig = model_configuration.get(app_id, {
            'thresholds': (0, np.inf)
        })
        self._low_threshold, self._high_threshold = appconfig['thresholds']
        logger.debug(repr(self) + ' Creating Model object')
        self._is_odd = lambda x: x % 2 == 1
        self._memory_is_full = lambda: len(self.memory) == self._memory_size
        logger.debug(
            repr(self) + ' Setting initial configuration of the model')
        self._memory_filled_flag = False
        logger.debug(repr(self) + ' Created')

    def _init_memory_fill_check(self):
        """Check if memory is filled, runs processing when memory fills up."""
        if self._memory_is_full() and np.nonzero(self.memory)[0].shape[0] > 0:
            logger.debug(
                repr(self) + 'Memory is fully populated, set flag to True')
            self._memory_filled_flag = True
            self._run_memory_calculations()

    def _validate_point(self, point):
        if isclose(
                point, 0, abs_tol=constants.CONVERT_NANO) or np.isnan(point):
            running_mean = self._get_running_mean_point()
            point = 0 if np.isnan(running_mean) else running_mean
        return point

    def _add_new_day(self, day_value):
        """Add new day value to memory."""
        logger.debug(
            repr(self) + 'Adding a new day to the memory {}'.format(day_value))
        validated_day_value = self._validate_point(day_value)
        self.memory.append(validated_day_value)
        logger.debug(repr(self) + 'A new day added.')
        if not self.memory_filled_flag:
            logger.debug(
                repr(self) + 'Memory set as unfull, checking memory fullness')
            self._init_memory_fill_check()
        else:
            logger.debug(
                repr(self) +
                'Calculating internal metrics for newest point in memory')
            self._run_single_new_point_calculations()

    def _check_state(self, val):
        """
        Check the input value against thresholds.

        Arguments :
        - val: int value to check against thresholds
        Returns:
        - int: -1 if below low threshold, 0 if between low and high,
                1 if above high
        """
        logger.debug(repr(self) + 'Checking state')
        if val > self.high_threshold:
            logger.debug(repr(self) + 'Value over high threshold, returning 1')
            return 1
        if val < self.low_threshold:
            logger.debug(
                repr(self) + 'Value below low threshold, returning -1')
            return -1
        logger.debug(repr(self) + 'Logger within thresholds, returning 0')
        return 0

    def _remove_zeros_from_memory(self):
        memory = np.asarray(self.memory)
        memory[memory == 0] = memory[np.min(np.nonzero(memory))]
        self.memory = deque(memory, maxlen=self._memory_size)

    def _run_memory_calculations(self):
        """Calculate internal fluctuation_metrics based on whole memory."""
        self._remove_zeros_from_memory()
        logger.debug(
            repr(self) + 'Calculating absolute internal fluctuation metrics')
        self.mean = self._get_running_mean_from_memory()
        logger.debug(repr(self) + 'Mean set to {}'.format(self.mean))
        self.fluctuation_metric = self._get_fluctuation_metric_from_memory()
        logger.debug(
            repr(self) +
            'Fluctuation metric set to {}'.format(self.fluctuation_metric))

    def _get_running_mean_from_memory(self):
        """Get running mean based on memory."""
        logger.debug(repr(self) + 'Getting running mean based on memory')
        result = deque(
            running_mean_fast(np.asarray(self.memory), self.window_size),
            maxlen=self.window_size * 2)
        logger.debug(
            repr(self) + 'Running mean based on memory: {}'.format(result))
        return result

    def _get_fluctuation_metric_from_memory(self):
        """Get fluctuation metric based on memory."""
        logger.debug(repr(self) + 'Getting fluctuation metric based on memory')
        result = deque(running_fn_norm(
                    np.asarray(self.memory), np.std,
                    self.window_size)[self.window_size:-self.window_size],
                maxlen=self.window_size)
        logger.debug(
            repr(self) +
            'Fluctuation metric based on memory: {}'.format(result))
        return result

    def _run_single_new_point_calculations(self):
        """Calculate internal metrics for newest point in memory."""
        logger.debug(
            repr(self) + 'Calculating internal metrics for newest point.')
        mean = self._get_running_mean_point()
        logger.debug(repr(self) + 'Calculated running mean: {}'.format(mean))
        self.mean.append(mean)
        metric = self._get_fluctuation_metric_point()
        logger.debug(
            repr(self) + 'Calculated fluctuation metric: {}'.format(metric))
        self.fluctuation_metric.append(metric)
        logger.debug(repr(self) + 'Newest mean and metric appended to memory.')

    def _get_metric_value(self, startpoint):
        preprocessed_signal = np.asarray(self.memory)[startpoint:
                                                      -self.half_window]
        std = np.std(preprocessed_signal)
        mean = self.mean[-self.half_window]
        return std / mean

    def _get_fluctuation_metric_point(self):
        """Get the value of newest fluctuation metric point."""
        logger.debug(
            repr(self) + 'Calculating fluctuation metric for newest point.')
        windowed_startpoint = -3 * self.half_window
        logger.debug(
            repr(self) + 'Windowed_startpoint: {}'.format(windowed_startpoint))
        result = self._get_metric_value(windowed_startpoint)
        logger.debug(repr(self) + 'Calculated metric: {}'.format(result))
        return result

    def _get_running_mean_point(self):
        """Get the value of newest running mean point."""
        result = np.mean(np.asarray(self.memory)[-self.window_size:])
        logger.debug(
            repr(self) +
            'Value of newest running mean point: {}'.format(result))
        return result

    @measure_time_metric
    def _run(self, incoming_signal):
        """
        Evaluate new datapoint in terms of fluctuations.

        Arguments:
        - incoming_signal: 1 day pd.dataframe containing fields
                           'value' and 'time'
        Returns:
        - dict: model output dict with pairs:
                'time_model': current time
                'value_model': current fluctuation state
                'app_id_model': appID
        """
        logger.debug(repr(self) + 'Starting evaluation')
        new_signal_value = incoming_signal.reset_index(drop=True).value[0]
        logger.debug(
            repr(self) + 'Value of latest signal: {}'.format(new_signal_value))
        current_time = incoming_signal.reset_index(drop=True).time[0]
        logger.debug(repr(self) + 'Current time: {}'.format(current_time))
        self._add_new_day(new_signal_value)
        if self.memory_filled_flag:
            current_state, mean_variance = self._predict()
            self._postprocess(current_state, mean_variance)
            model_output = {
                'time_model': current_time,
                'value_model': current_state,
                'app_id_model': self.app_id,
            }
        else:
            model_output = {
                'time_model': current_time,
                'value_model': None,
                'app_id_model': self.app_id,
            }
        logger.debug(repr(self) + 'Model output dict: {}'.format(model_output))
        return model_output

    def _predict(self):
        """Return model state and current value upon predition."""
        logger.debug(repr(self) + 'Enough data for evaluation, proceeding')
        logger.debug(repr(self) + 'Calculate current state')
        mean_variance = np.mean(
            np.asarray(self.fluctuation_metric)[-self.window_size:])
        current_state = self._check_state(mean_variance)
        logger.debug(repr(self) + 'Output of Model'.format(current_state))
        return current_state, mean_variance

    def _postprocess(self, current_state, mean_variance):
        """Push current state and value to prometheus."""
        logger.debug(repr(self) + 'Send current_state to prometheus')
        model_output_prometheus.labels(self.app_id,
                                       self.metric).set(current_state)
        model_metric_prometheus.labels(self.app_id,
                                       self.metric).set(mean_variance)

        logger.debug(repr(self) + 'current_state sent to prometheus')

    @property
    def window_size(self):
        return self._window_size

    @property
    def half_window(self):
        return self._half_window

    @property
    def memory_size(self):
        return self._memory_size

    @property
    def low_threshold(self):
        return self._low_threshold

    @property
    def high_threshold(self):
        return self._high_threshold

    @property
    def memory_filled_flag(self):
        return self._memory_filled_flag

    def __repr__(self):
        """Use repr(self) for DEBUG, ERROR, CRITICAL logging levels."""
        return '{} - {} - Low_thresh {} - High_thresh {}: '.format(
            self.app_id, self.metric, self.low_threshold, self.high_threshold)

    def __str__(self):
        """Use str(self) for INFO, WARNING logging levels."""
        return '{}: '.format(self.app_id)


class PeriodicityRemovalMidtermFluctuationDetector(MidTermFluctuationDetector):
    def _get_metric_value(self, startpoint):
        preprocessed_signal = (self._get_memory_without_periodicity(
            PERIOD, startpoint) + np.mean(
                np.asarray(
                    self.memory)[startpoint:-self.half_window]))
        mean = np.abs(np.mean(preprocessed_signal))
        std = np.std(preprocessed_signal)
        return std / mean

    def _get_memory_without_periodicity(self, lag, windowed_startpoint):
        return (np.asarray(
            self.memory)[windowed_startpoint:-self.half_window] - np.asarray(
                self.memory)[windowed_startpoint - lag:-self.half_window - lag]
                )

    def _get_fluctuation_metric_from_memory(self):
        """Get fluctuation metric based on memory."""
        logger.debug(repr(self) + 'Getting fluctuation metric based on memory')
        preprocessed_signal = np.asarray(self.memory, dtype='float64')
        preprocessed_signal[PERIOD:] -= preprocessed_signal[:-PERIOD]
        preprocessed_signal[PERIOD:] += np.mean(np.asarray(self.memory))
        result = deque(
            running_fn_norm(preprocessed_signal, np.std, self.window_size,
                            False)[self.window_size:-self.window_size],
            maxlen=self.window_size)
        logger.debug(
            repr(self) +
            'Fluctuation metric based on memory: {}'.format(result))
        return result
