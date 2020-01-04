from src.decisionmakers import DecisionMaker
from prometheus_client import Gauge
from src.utils import measure_time_metric, get_time_dict
import numpy as np
import datetime
import logging

logger = logging.getLogger('root')

model_output_prometheus = Gauge('Decisionmaker_output', 'appIDs',
                                ['app_id', 'metric_name'])


class FluctuationDecisionMaker(DecisionMaker):
    """
    Component deciding what message to send based on fluctuation model.

    FluctuationDecisionMaker decides whether to send notifications that the
    fluctuations have become abnormal, or returned to normal range based on
    model output
    """

    def __init__(self,
                 app_id=None,
                 metric='number of calls',
                 aid_service_connection=None,
                 message_lag=21):
        """
        Create decisionmaker object initialized at stable state.

        Arguments:
        - app_id: integer
        - metric: string specifying the metric
        - aid_service_connection: string
        """
        super(FluctuationDecisionMaker,
              self).__init__(app_id, metric, aid_service_connection, self._run)
        logger.debug(repr(self) + ' Creating FluctuationDecisionMaker object ')
        self.message_lag = message_lag
        self.current_state = None
        self.destabilized_dt = None
        self.messages = {
            0: 'The signal fluctuations have returned back to normal',
            1: 'The signal fluctuations have become abnormally high',
            -1: 'The signal fluctuations have become abnormally low'
        }
        self._init_grpc_messages()
        logger.debug(repr(self) + ' FluctuationDecisionMaker created')

    def _pick_message(self, state):
        """
        Pick message based on current and new state.

        Arguments:
        - state: -1, 0 or 1, new state returned by the model

        Returns:
        - result: string with appropriate message

        """
        logger.debug(repr(self) + 'Picking message for provided state.')
        result = self.messages[state]
        logger.debug(repr(self) + 'Message picked {}'.format(result))
        return result

    def _send_grpc_message(self, state, time):
        """
        Send grpc message based on current and new state.

        Arguments:
        - state: -1, 0 or 1, new state returned by the model
        - time: time extracted from last date_package
        """
        if None not in (self.grpc_messages, state):
            logger.debug(
                repr(self) + 'Sending grpc message state {}'.format(state))
            if state != 0:
                self.destabilized_dt = (time - datetime.timedelta(
                                                days=self.message_lag))
                ts = get_time_dict(self.destabilized_dt, current=(0, None))
            elif state == 0:
                ts = get_time_dict(self.destabilized_dt, current=(
                                    0, (time -
                                        self.destabilized_dt -
                                        datetime.timedelta(
                                            days=self.message_lag)).days))
            self.grpc_messages[state](time, ts, self.app_id)
            logger.info(repr(self) + 'Grpc message sent {}.'.format(
                                                self.grpc_messages[state]))

    @measure_time_metric
    def _run(self, model_output):
        """
        Evaluate the fluctuation state change and appropriate message.

        Arguments:
        - model_output: a dictionary with model output dict with pairs:
                        'time_model' : current time
                        'value_model' : current fluctuation state
                        'app_id_model' : appID

        Returns:
        - output: decisionmaker output dict with pairs:
                  'time_model': current time
                  'value_model': current fluctuation state
                  'app_id_model': appID
                  'message': appropriate message
        """
        logger.debug(repr(self) + 'Starting _run method.')
        model_verdict, current_time = self._preprocess(model_output)
        output = self._decide(model_verdict, current_time)
        self._postprocess(model_verdict)
        return output

    def _preprocess(self, model_output):
        """

        Arguments:
        - model_output: a dictionary with model output dict with pairs:
                        'time_model' : current time
                        'value_model' : current fluctuation state
                        'app_id_model' : appID

        Returns:
        - model_verdict: current fluctuation state
        - current_time: time extracted from last date_package

        """
        current_time, model_verdict = self._unpack_output(model_output)
        self._destabilize_dt(current_time)
        self._check_current_state(model_verdict)
        return model_verdict, current_time

    def _decide(self, model_verdict, current_time):
        """
        Arguments:
        - model_verdict: current fluctuation state
        - current_time: time extracted from last date_package

        Returns:
        - output: decisionmaker output dict with pairs:
                  'time_model': current time
                  'value_model': current fluctuation state
                  'app_id_model': appID
                  'message': appropriate message

        """
        message = self._get_message(model_verdict, current_time)
        self._update_state(model_verdict)
        output = self._set_output(current_time, model_verdict, message)
        return output

    def _postprocess(self, verdict):
        """
        Update prometheus variables.

        Arguments:
        - verdict: current fluctuation state
        """
        logger.debug(repr(self) + 'Start post-processing of the output')
        prometheus_input = verdict if verdict is not None else np.nan
        logger.debug(
            repr(self) + 'Current State: {}.'.format(self.current_state))
        model_output_prometheus.labels(self.app_id,
                                       self.metric).set(prometheus_input)
        logger.debug(
            repr(self) +
            'State sent to prometheus {}.'.format(prometheus_input))

    def _unpack_output(self, output):
        """
        Arguments:
        - output: decisionmaker output dict with pairs:
                  'time_model': current time
                  'value_model': current fluctuation state
                  'app_id_model': appID
                  'message': appropriate message

        Returns:
        - current_time: time extracted from last date_package
        - model_verdict: current fluctuation state

        """
        logger.debug(repr(self) + 'Model output: {}'.format(output))
        current_time = output['time_model']
        model_verdict = output['value_model']
        return current_time, model_verdict

    def _destabilize_dt(self, time):
        """Destabilize time passed as an argument."""
        if self.destabilized_dt is None:
            self.destabilized_dt = time

    def _check_current_state(self, verdict):
        """Adjust state if verdict exists."""
        if self.current_state is None and verdict == 0:
            self.current_state = 0

    def _get_message(self, verdict, current_time):
        """
        Get currently appropriate message.

        Arguments:
        - verdict: current fluctuation state
        - current_time: time extracted from last date_package

        Returns:
        - message: string with appropriate message or None

        """
        if verdict != self.current_state and verdict is not None:
            logger.debug(
                repr(self) + 'State has changed, picking up new message.')
            message = self._pick_message(verdict)
            logger.debug(repr(self) + 'Message set to {}.'.format(message))
            self._send_grpc_message(verdict, current_time)
            return message

    def _update_state(self, verdict):
        """
        Arguments:
        - verdict: current fluctuation state
        """
        logger.debug(
            repr(self) + 'Setting Current State of the Decision Maker.')
        self.current_state = verdict
        logger.debug(
            repr(self) + 'Current State: {}.'.format(self.current_state))

    def _set_output(self, current_time, model_verdict, message):
        """
        Arguments:
        - current_time: time extracted from last date_package
        - model_verdict: current fluctuation state
        - message: appropriate message

        Returns:
        - output: dict with pairs:
                  'time_model': current time
                  'value_model': current fluctuation state
                  'app_id_model': appID
                  'message': appropriate message

        """
        logger.debug(repr(self) + 'Setting Decision Maker Output.')
        output = {
            'time_model': current_time,
            'value_model': model_verdict,
            'app_id_model': self.app_id,
            'message': message
        }
        logger.debug(
            str(self) + 'Decision Maker Output for {} metric: {}'.format(
                self.metric, output))
        return output

    def _init_grpc_messages(self):
        """Create message dictionary."""
        if self.aid_service_connection:
            self.grpc_messages = self.aid_service_connection.messages[
                self.metric]['fluctuation']
        else:
            self.grpc_messages = None

    def __repr__(self):
        """Use repr(self) for DEBUG, ERROR, CRITICAL logging levels."""
        return '{} - Metrics {}: '.format(self.app_id, self.metric)

    def __str__(self):
        """Use str(self) for INFO, WARNING logging levels."""
        return '{}: '.format(self.app_id)
