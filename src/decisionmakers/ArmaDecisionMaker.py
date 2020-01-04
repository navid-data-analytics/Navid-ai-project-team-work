"""Class for Arma decision making component."""
from src.decisionmakers import DecisionMaker
from src.utils import measure_time_metric, get_time_dict
from prometheus_client import Gauge
import numpy as np
import logging

output_prometheus = Gauge('ArmaDecisionMaker_output', 'appIDs',
                          ['app_id', 'metric_name'])
rate_of_change_prometheus = Gauge('Prediction_rate__of_change', 'appIDs',
                                  ['app_id', 'metric_name'])
logger = logging.getLogger('root')


class ArmaDecisionMaker(DecisionMaker):
    """
    A ArmaDecisionMaker object decides what should be sent as a verdict.

    The component is a DecisionMaker that takes model output as input and
    decides what is the message sent forward. Additionally it sends predicted
    signal for Grafana.
    """

    def __init__(self,
                 app_id=None,
                 metric='number of calls',
                 aid_service_connection=None,
                 predict_interval=30):
        """
        Construct a new ArmaDecisionMaker instance.

        Arguments:
        - app_id: integer
        - metric: string specifying the metric
        - aid_service_connection: string
        - predict_interval: int, number of days to predict forward
        """
        super(ArmaDecisionMaker,
              self).__init__(app_id, metric, aid_service_connection, self._run)
        self._predict_interval = predict_interval
        logger.debug(repr(self) + ' Creating ArmaDecisionMaker object ')
        self.message = None
        self.decision_decoder = {
            1: 'grow',
            -1: 'decline',
        }
        self._init_grpc_messages()
        logger.debug(repr(self) + ' ArmaDecisionMaker created')

    def _init_grpc_messages(self):
        """Create message dictionary."""
        if self.aid_service_connection:
            self.grpc_messages = self.aid_service_connection.messages[
                self.metric]['prediction']
        else:
            self.grpc_messages = None

    @measure_time_metric
    def _run(self, model_output):
        """
        Run is called each time a dataframe is pushed into the component.

        Arguments:
        - model_output: a dictionary containing prediction, rate of change
                        and current time field.
        Returns:
        - output: a dictionary containing time_model, app id model,
                  decision, rate of change fields.
        """
        prediction, current_time = self._preprocess(model_output)
        decision = self._decide()
        output = self._post_process(decision, prediction, current_time)
        logger.debug(
            repr(self) +
            'Output of the Prediction Decision Maker {}'.format(output))
        return output

    def _preprocess(self, model_output):
        """
        Preprocess the signal.

        Returns:
        - prediction: dataframe with predicted values by the model
        - current_time: time extracted from last date_package
        """
        prediction, rate_of_change, current_time = model_output
        logger.debug(
            repr(self) + 'prediction {} rate_of_change {} current_time {}'.
            format(prediction, rate_of_change, current_time))
        self.rate_of_change = rate_of_change
        return prediction, current_time

    def _decide(self):
        """
        Decide whether signal is predicted to grow or decline.

        The decision is based on rate of change passed by ArmaPredictor,
        where percentage difference is calculated.

        Returns:
        - decision: 1, 0, -1 or np.nan depending on the scenario
        """
        decision = self._handle_decision()
        self._set_message(decision)
        return decision

    def _handle_decision(self):
        """
        Handle if decision is possible.

        If there is not enough data, set the message to 'Not enough data
        for evalution'.

        Returns:
        - decision: 1, 0, -1 or np.nan depending on the scenario
        """
        if self.rate_of_change is not None:
            decision = 1 if self.rate_of_change >= 0 else -1
            logger.debug(repr(self) + 'Decision assigned: {}'.format(decision))
        else:
            self.message = 'Not enough data for evaluation.'
            logger.debug(
                repr(self) + 'Rate of change is None, message set to {}'.
                format(self.message))
            decision = np.nan
        return decision

    def _post_process(self, decision, prediction, current_time):
        """
        Prepare the message and Prometheus logs.

        Both decision and prediction are processed for the message,
        additionally passing the forecasting signal for grafana.

        Arguments:
        - decision: 1, 0, -1 or np.nan depending on the scenario
        - prediction: dataframe containing predicted 'time' and 'value' columns
        - current_time: time extracted from date_package

        Returns:
        - output: a dictionary containing time_model, app id model,
                  decision, rate of change fields.
        """
        logger.debug(repr(self) + 'Start post-processing of the output')
        self._update_prometheus(decision)
        self._try_sending_grpc(current_time)
        output = self._get_output(current_time, decision)
        return output

    def boundary_check(self, time):
        raise NotImplementedError('This method is not imp')

    def _get_output(self, time, decision):
        """
        Create output dictionary.

        Arguments:
        - time: time extracted from date_package
        - decision: 1, 0, -1 or np.nan depending on the scenario

        Returns:
        - output: a dictionary containing time_model, app id model,
                  decision, rate of change fields.
        """
        output = {
            'time_model': time,
            'app_id_model': self.app_id,
            'decision': decision,
            'rate_of_change': self.rate_of_change
        }
        logger.debug(
            repr(self) + 'Output of the ArmaDecisionMaker'.format(output))
        return output

    def _set_message(self, decision):
        """
        Set message if empty.

        Arguments:
        - decision: 1, 0, -1 or np.nan depending on the scenario
        """
        if self.message is None:
            logger.debug(repr(self) + 'Message is not set, setting message...')
            self.message = 'Your {} is expected to {}'.format(
                self.metric, self.decision_decoder[decision]) + \
                ' by {:.2f}% in the next {} days.'.format(
                self.rate_of_change, self.predict_interval)
            logger.debug(repr(self) + 'Message set to {}'.format(self.message))

    def _update_prometheus(self, decision):
        """
        Update Prometheus with decision and rate_of_change.

        Arguments:
        decision: 1, 0, -1 or np.nan depending on the scenario
        """
        if decision is None:
            output_prometheus.labels(self.app_id, self.metric).set(np.nan)
        else:
            output_prometheus.labels(self.app_id, self.metric).set(decision)
        logger.debug(
            repr(self) + 'Decision sent to prometheus {}'.format(decision))
        arma_rate_of_change = self.rate_of_change if\
            self.rate_of_change is not None else 0
        logger.debug(
            repr(self) +
            'rate_of_change sent to prometheus {}'.format(arma_rate_of_change))
        rate_of_change_prometheus.labels(self.app_id,
                                         self.metric).set(arma_rate_of_change)

    def _try_sending_grpc(self, time):
        """
        Check AID-E connection and send message.

        Continue sending process if AID-E is connected,
        otherwise send error to Sentry

        Arguments:
        - time: time extracted from last date_package.
        """
        if self.aid_service_connection:
            if time is not None:
                logger.debug(
                  repr(self) + 'Connection to AID-E allowed, send themessage.')
                periods = self._prepare_message(time)
                if periods is not None:
                    self._send_message(time, periods)
        else:
            logger.error(repr(self) + 'No connection to AID-E!')

    def _prepare_percentage(self):
        """
        Prepare percentage from self.rate_of_change.

        Returns:
        - percentage
        """
        logger.debug(
            repr(self) + 'Rate of change not None, calculate percentage')
        percentage = round(self.rate_of_change)
        logger.debug(
            repr(self) + 'Calculated percentage: {}'.format(percentage))
        return percentage

    def _prepare_periods(self, time):
        """
        Create dictionary with accordance to documentation.

        Arguments:
        - time: Time in DateTime format

        Returns:
        - periods: Dictionary of time periods with accordance to documentation
        """
        logger.debug(repr(self) + 'Creating Timestamp dictionary')
        periods = get_time_dict(time,
                                current=(-self.predict_interval, 0),
                                future=(0, self.predict_interval))
        logger.debug(
            repr(self) + 'Timestamp dictionary done: {}'.format(periods))
        return periods

    def _prepare_message(self, time):
        """
        Decide which message to send to AID-E.

        Arguments:
        - time: Time in DateTime format

        Returns:
        - periods: Dictionary of time periods with accordance to documentation
        """
        should_prepare = self.boundary_check(time)
        if not should_prepare:
            return
        if self.rate_of_change is None:
            logger.debug(
                repr(self) + 'Rate of change is None,'
                'no prediction sent yet.')
            return
        elif np.isnan(np.array([self.rate_of_change], dtype=np.float64)):
            logger.error(
                repr(self) +
                '{}: Calculated rate of change is NaN!'.format(time))
            logger.error(repr(self) + 'AppID: {}'.format(self.app_id))
            logger.error(repr(self) + 'Cannot send message to Grpc!')
            return
        periods = self._prepare_periods(time)
        return periods

    @property
    def predict_interval(self):
        """The interval we make decision upon."""
        return self._predict_interval

    def __repr__(self):
        """Use repr(self) for DEBUG, ERROR, CRITICAL logging levels."""
        return '{} - Metrics {} - Predict_interval {}: '.format(
            self.app_id, self.metric, self.predict_interval)

    def __str__(self):
        """Use str(self) for INFO, WARNING logging levels."""
        return '{}: '.format(self.app_id)
