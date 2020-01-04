from src.decisionmakers import TrendDecisionMaker
import logging
import numpy as np
from src.utils import measure_time_metric, get_time_dict
from prometheus_client import Gauge

dm_immediate_trend_decision_prometheus = Gauge(
    'dm_immediate_trend_decision_short', 'appIDs', ['app_id', 'metric_name'])
dm_immediate_change_to_report_prometheus = Gauge(
    'dm_immediate_change_to_report_short', 'appIDs', ['app_id', 'metric_name'])

logger = logging.getLogger('root')


class DecisionMakerShortTerm(TrendDecisionMaker):
    """
    This class is responsible for deciding what kind
    and when a notification should be created.
    """

    def __init__(self,
                 aid_service_connection=None,
                 app_id=None,
                 metric=None,
                 report_days=(None, ),
                 reportable_change_bounds=(0, 1000),
                 days_to_report=14):
        """
        Construct a new DecisionMakerMidTerm instance.

        Arguments:
        - aid_service_connection: MessageClient object, if None the component
                                  will not send grpc messages
        - app_id: integer
        - metric: string specifying the metric
        - report_days: tuple of ints, which days to send monthly messages at
        - reportable_change_bounds: tuple of ints, what are the absolute
                                    boundaries of a valid change percentage
        - days_to_report: how many days are included within a report
        """
        super(DecisionMakerShortTerm, self).__init__(
            aid_service_connection, app_id, metric, report_days,
            reportable_change_bounds, days_to_report)
        self._init_grpc_messages()

    def _init_grpc_messages(self):
        """Create message dictionary."""
        if self.aid_service_connection:
            logger.debug(self.aid_service_connection.messages)
            self._grpc_messages = self.aid_service_connection.messages[
                self.metric]['trend']

            self._message_types = {
                'immediate': {
                    1: 'Immediate_up_short',
                    -1: 'Immediate_down_short'
                }
            }

    @measure_time_metric
    def _run(self, model_output):
        """
        Check if a monthly and immediate report should be sent.

        Arguments:
        - model_output: dict, the input of the Class coming from Detector

        Returns:
        - decision_maker_output: a dict with following structure:
                                 'time_model': current_time,
                                 'value_model': model_decision,
                                 'app_id_model': app_id,
                                 'app_status': app_status,
                                 'decision': trend_decision
        """
        logger.debug(
            repr(self) +
            'Input of the Short Trend Decision Maker {}'.format(model_output))
        model_decision, no_traffic_flag, current_time = \
            self._preprocess(model_output)
        decision_maker_output, immediate_decision, immediate_change = \
            self._decide(current_time, model_decision, no_traffic_flag)
        self._postprocess(immediate_decision, immediate_change)
        logger.debug(
            repr(self) + 'Output of the Short Trend Decision Maker {}'.format(
                decision_maker_output))
        return decision_maker_output

    def _decide(self, current_time, model_decision, no_traffic_flag):
        """
        Make decision and store it in output dictionary.

        Arguments:
        - model_decision: int, is an integer which 1 for growth, -1 for decline
                          and 0 for no trend
        - no_traffic_flag: bool, if there have been 3 days with no traffic
        - current_time: time extracted from last date_package

        Returns:
        - decision_maker_output: a dict with following structure:
                                 'time_model': current_time,
                                 'value_model': model_decision,
                                 'app_id_model': app_id,
                                 'app_status': app_status,
                                 'decision': trend_decision
        - change_to_report: float, the amount of growth to report
        - monthly_decision: int, is 1 if there has been some growth in
                            the past 15 days to be reported. -1 for declines
                            and zero for nothing
        - trend_decision: is an integer which 1 for growth, -1 for decline
                          and 0 for no trend
        """
        logger.debug(repr(self) + 'in decision function.')
        immediate_decision, immediate_change = self._handle_immediate_check(
            no_traffic_flag, model_decision, current_time)
        logger.debug(
            repr(self) + 'immediate decision {} immediate change {}'.format(
                immediate_decision, immediate_change))
        decision_maker_output = {
            'time_model': current_time,
            'value_model': model_decision,
            'app_id_model': self.app_id,
            'send': self._send,
            'trend_monthly_decision': None,
            'trend_immediate_decision': immediate_decision,
            'message_details': self._message_details
        }

        logger.debug(
            repr(self) +
            'short_decision_maker output: {}'.format(decision_maker_output))
        return decision_maker_output, immediate_decision, immediate_change

    def _postprocess(self, immediate_decision, immediate_change):
        """Update Prometheus with decision and rate_of_change.

        Arguments:
        - immediate_decision: is an integer which 1 for growth, -1 for decline
                          and 0 for no trend
        - immediate_change: float, the amount of change to report
        """
        dm_immediate_trend_decision_prometheus.labels(
            self.app_id, self.metric).set(immediate_decision)
        dm_immediate_change_to_report_prometheus.labels(
            self.app_id, self.metric).set(immediate_change)
