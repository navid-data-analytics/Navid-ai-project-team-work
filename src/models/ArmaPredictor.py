"""ArmaPredictor class implementing Autoregressive Moving Average model."""
from src.models import Predictor
from src.utils.measures import measure_time_metric
from src.utils.running_functions import running_mean_fast
from statsmodels.tsa.seasonal import seasonal_decompose
from statsmodels.tsa.arima_model import ARMA
from prometheus_client import Gauge
import src.utils.constants as constants
import pandas as pd
import numpy as np
import logging
import math

output_prometheus = Gauge('ArmaPredictor_output', 'appIDs',
                          ['app_id', 'metric_name'])

logger = logging.getLogger('root')


class ArmaPredictor(Predictor):
    """
    An ArmaPredictor object is an instance of ARMA forecasting model.

    The component is a predictor that takes dataframe on the input,
    and once a certain history threshold is filled it starts to predict
    future values. The threshold is set in model_configuration dictionary
    with 'history_size' integer value.
    """

    @staticmethod
    def add_whitenoise(history):
        logger.debug('Adding whitenoise to the signal.')
        whitenoise = np.random.normal(
            loc=0, scale=0.1, size=(history.shape[0], 1))
        abs_of_whitenoise = np.abs(whitenoise)
        history += abs_of_whitenoise
        logger.debug('Whitenoise added to the signal.')
        return history

    def __init__(self,
                 model_configuration,
                 app_id=None,
                 predict_interval=7,
                 metric=None):
        """
        Construct a new Predictor instance.

        Arguments:
        - model_configuration: dictionary containing model_configuration
            where key 'history_size' defines model capacity and 'arma_order'
            defines dictionary with appid to parameter mapping
        - app_id: integer, self-explainatory
        - predict_interval: integer, number of days taken by the model while
            training

        """
        self._app_id = app_id
        self._predict_interval = predict_interval
        self._interval_index = predict_interval
        self._order = model_configuration['arma_order'].get(app_id, (1, 0))
        self._window_size = model_configuration['seasonality_period']
        logger.debug(
            repr(self) + ' Creating ArmaPredictor object {}'.format(app_id))
        Predictor.__init__(self, model_configuration, app_id, metric)
        logger.debug(
            repr(self) + ' Setting initial configuration of the model')
        self._is_trained = False
        self.build_model = ARMA
        self.not_empty = lambda x: np.sum(x[-5:]) != 0
        logger.debug(repr(self) + ' ArmaPredictor {} created'.format(app_id))

    @measure_time_metric
    def _run(self, incoming_df):
        """
        Run is called each time a dataframe is pushed into the component.

        Arguments:
        - incoming_df: pandas dataframe containing 'time' and 'value' columns

        Returns:
        - result: a tuple containing:
                  prediction - a forecast  if predicted otherwise None
                  rate_of_change - float, predicted change in slope
                                  given as percentage
                  current_time - time extracted from last date_package

        """
        current_time = self._preprocess(incoming_df)
        self._update_model(incoming_df)
        if self.is_trained:
            logger.debug(repr(self) + 'Model is trained, starting prediction.')
            logger.debug(repr(self) + 'current_time: {}'.format(current_time))
            prediction = self._predict()
            logger.debug(repr(self) + 'prediction: {}'.format(prediction))
            rate_of_change = self._postprocess(prediction, current_time)
            logger.debug(
                repr(self) + 'rate_of_change: {}'.format(rate_of_change))
            result = (prediction, rate_of_change, current_time)
        else:
            logger.debug(
                repr(self) + 'Model is not trained, sending mock result.')
            result = (None, None, current_time)
            logger.debug(repr(self) + str(result))
        return result

    def _update_model(self, incoming_df):
        """
        Update history and fit ARMA to the series if possible.

        Arguments:
        - incoming_df: pandas dataframe containing 'time' and 'value' columns
        """
        logger.debug(
            repr(self) + 'Updating history with {}'.format(incoming_df))
        self.history.append(incoming_df)
        if self.filled(self.history):
            logger.debug(repr(self) + 'Model history is full, proceeding.')
            self.preprocessed_history = self._preprocess_model_history()
            if self.not_empty(self.preprocessed_history):
                logger.debug(repr(self) + 'Data is not empty, proceeding.')
                self.model = self._prepare_model()
                self._is_trained = True
                logger.debug(repr(self) + 'Set is_trained flag to True.')

    def _preprocess_model_history(self):
        """
        Preprocess model history in a way it is feedable to ARMA model.

        Returns:
        - preprocessed_history: np.array with time and value for each datapoint
        """
        logger.debug(repr(self) + 'Start preprocessing.')
        history_df = pd.concat(self.history)
        logger.debug(repr(self) + 'History deque has been concatenated.')
        history_df = history_df[['time', 'value']]
        history_df.set_index('time', inplace=True)
        preprocessed_history = np.asarray(history_df).astype(np.float)
        logger.debug(repr(self) + 'History converted to numpy array.')
        self._raw_history = self.add_whitenoise(preprocessed_history)
        decomposed_signal = seasonal_decompose(
            list(preprocessed_history), freq=self._window_size)
        half_window = math.floor(self._window_size / 2)
        assert half_window == 3, half_window
        assert isinstance(half_window, int), type(half_window)
        preprocessed_signal = decomposed_signal.trend[half_window:-half_window]
        preprocessed_signal = running_mean_fast(preprocessed_signal,
                                                self._window_size)
        logger.debug(
            repr(self) +
            'History decomposed to {}'.format(preprocessed_signal))
        logger.debug(
            repr(self) + 'Preprocessed signal shape {} (should be 150!)'.
            format(preprocessed_signal.shape))
        logger.debug(
            repr(self) +
            'Returning preprocessed signal: {}'.format(preprocessed_signal))
        return preprocessed_signal

    def _prepare_model(self):
        """
        Prepare model for forecast.

        Arguments:
        - preprocessed_history: np.array with time and value for each datapoint

        Returns:
        - model: ARMA model fit to provided history
        """
        try:
            logger.debug(
                repr(self) + 'Fitting preprocessed_history to the model.')
            model = self.build_model(self.preprocessed_history, self.order)
            model = model.fit(disp=0)
            logger.debug(repr(self) + 'Model fit correctly.')
        except Exception:
            logger.warning(
                repr(self) +
                'The model with whitenoise did not fit, use (1,0)')
            model = self.build_model(self._raw_history, (1, 0))
            model = model.fit(disp=0)
            logger.debug(repr(self) + 'Model fit correctly.')
        return model

    def _predict(self):
        """
        Forecast next X days, X is defined by self.predict_interval field.

        Returns:
        - forecast: pandas dataframe containing 'time' and 'value' columns
        """
        logger.debug(repr(self) + 'Starting prediction.')
        forecast, _, __ = self.model.forecast(steps=self.predict_interval)
        logger.debug(repr(self) + 'Success, prediction {}'.format(forecast))
        return forecast

    def _postprocess(self, prediction, current_time):
        """
        Send relevant log to prometheus.

        Arguments:
        - prediction: forecast, if predicted otherwise None

        Returns:
        - rate: float, predicted change in slope given as percentage
        """
        logger.debug(repr(self) + 'Start post-processing of the output')
        last_forecasted_value = np.mean(prediction)
        logger.debug(
            repr(self) +
            'Last forecasted value: {}'.format(last_forecasted_value))
        last_known_value = np.mean(
            self.preprocessed_history[self.interval_index:])
        logger.debug(
            repr(self) + 'Last known value: {}'.format(last_known_value))
        rate = self._calculate_change_rate(last_known_value,
                                           last_forecasted_value)
        logger.debug(repr(self) + 'Calculated rate: {}'.format(rate))
        rate = abs(rate)
        logger.debug(repr(self) + 'Absolute value of rate: {}'.format(rate))
        output_prometheus.labels(self.app_id, self.metric).set(rate)
        logger.debug(repr(self) + 'Rate sent to prometheus {}'.format(rate))
        output = {
            'time_model': current_time,
            'rate_of_change': rate,
            'app_id_model': self.app_id
        }
        logger.debug(
            repr(self) + 'Output of the ArmaPredictor: {}'.format(output))
        return rate

    def _calculate_change_rate(self, known_value, forecasted_value):
        """
        Calculate rate of change for known and forecasted signal.

        Arguments:
        - known_value: float, last value for history timeseries
        - forecasted_value: float, last value for forecasted signal

        Returns:
        result: float, rate of change given as percentage (between 0. - 100.)
        """
        logger.debug(
            repr(self) + 'If last known value is 0, substitute  it to 1')
        known_value = 1 if known_value == 0 else known_value
        logger.debug(
            str(self) + 'Forecasted value {}'.format(forecasted_value))
        logger.debug(str(self) + 'Last known value {}'.format(known_value))
        difference = forecasted_value - known_value
        result = float((difference) * constants.CONVERT_PERCENT / known_value)
        logger.debug(
            repr(self) + 'Percentage rate of change {}'.format(result))
        return result

    def _preprocess(self, incoming_df):
        """Extract time from the input signal."""
        logger.debug(repr(self) + 'Starting _run method.')
        current_time = incoming_df.reset_index(drop=True).time[0]
        logger.debug(repr(self) + 'current_time: {}'.format(current_time))
        return current_time

    @property
    def is_trained(self):
        """Boolean flag showing if an object was trained."""
        return self._is_trained

    @property
    def predict_interval(self):
        """The interval we make decision upon."""
        return self._predict_interval

    @property
    def interval_index(self):
        return self._predict_interval

    @property
    def order(self):
        return self._order

    @property
    def window_size(self):
        return self._window_size

    def __repr__(self):
        """Use repr(self) for DEBUG, ERROR, CRITICAL logging levels."""
        return '{} - Order {} - Predict_interval {}: '.format(  #noqa
            self.app_id, self.order, self.predict_interval)

    def __str__(self):
        """Use str(self) for INFO, WARNING logging levels."""
        return '{}: '.format(self.app_id)
