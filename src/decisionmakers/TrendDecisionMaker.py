from src.decisionmakers import DecisionMaker
import logging
import numpy as np
from src.utils import measure_time_metric, get_time_dict

logger = logging.getLogger('root')


class TrendDecisionMaker(DecisionMaker):
    """
    This class is responsible for deciding what kind
    and when a notification should be created.
    """

    def __init__(self,
                 aid_service_connection=None,
                 app_id=None,
                 metric=None,
                 report_days=(1, 16),
                 reportable_change_bounds=(0, 1000),
                 days_to_report=60):
        """
        Construct a new TrendDecisionMaker instance.

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
        super(TrendDecisionMaker,
              self).__init__(app_id, metric, aid_service_connection, self._run)
        self._days_to_report = days_to_report
        self._report_days = report_days
        self._reportable_change_bounds = reportable_change_bounds
        self._init_grpc_messages()
        self._init_trend_monitor_vars()

    def _init_grpc_messages(self):
        self._grpc_messages = self.aid_service_connection.messages[
            self.metric]['trend']

        self._message_types = {
            'immediate': {
                1: 'Immediate_up',
                -1: 'Immediate_down'
            },
            'monthly': {
                1: '15_days_up',
                -1: '15_days_down'
            }
        }

    def _run(self, model_output):
        raise NotImplementedError("Method not implemented!")

    def _decide(self, current_time, model_decision, no_traffic_flag):
        raise NotImplementedError("Method not implemented!")

    def _postprocess(self, immediate_decision, immediate_change,
                     monthly_decision, monthly_change):
        raise NotImplementedError("Method not implemented!")

    def _init_trend_monitor_vars(self):
        self._trends_list = []
        self._current_trend_direction = 0

    def _get_monthly_trend_direction(self):
        """
        Check if there has been only downward trend in the past month.

        Returns:
        - True for only downward trend, False otherwise
        """
        if len(self._trends_list) > 0:
            return int(sum(self._trends_list) / len(self._trends_list))
        return 0

    @staticmethod
    def get_change_metric(new_val, old_val):
        """
        Calculate the growth or decline percentage.

        Arguments:
        - new_val: float, new value for measured metric
        - old_val: float, old value for measured metric

        Returns:
         - change: float, amount of change percentage
        """
        change = (((new_val - old_val) / max(
            abs(old_val),
            np.finfo(float).eps) * np.sign(old_val) * 100))
        change_rounded = round(change)
        return change_rounded

    def _is_valid_change(self, change, direction):
        """
        Check if the change direction and magnitude are valid.

        Arguments:
        - change: int, change value to be validated
        - direction: 1 or -1, for positive and negative slope

        Returns:
        - True if valid, False otherwise
        """
        return (self.reportable_change_bounds[0] < change * direction <
                self.reportable_change_bounds[1])

    def _prepare_immediate_message_details(self, ts, scores,
                                           change_percentage):
        return [
            self._datetime_utc, ts, scores, self.app_id,
            int(np.abs(change_percentage)), self._days_to_report
        ]

    def _prepare_monthly_message_details(self, ts, scores, change_percentage):
        return [
            self._datetime_utc, ts, scores, self.app_id,
            int(np.abs(change_percentage))
        ]

    def _send_message_aid_service_connection(self, change_percentage,
                                             change_direction, mtype):
        """
        Check if the aid_service_connection has been initiated, write to db.

        Arguments:
        - change_percentage: float, the amount of change to report
        - change_direction: int, 1 or -1 for a upwards or downwards slope
        - mtype: string, either 'monthly' or 'immediate'
        """
        assert mtype in ('monthly', 'immediate')

        ts = get_time_dict(
            self._datetime_utc,
            previous=(-self.days_to_report, -self.days_to_report/2),
            current=(-self.days_to_report/2, 0))
        scores = (self._avg_traffic_old, self._avg_traffic_new)

        prepare_md = (self._prepare_monthly_message_details
                      if mtype == 'monthly' else
                      self._prepare_immediate_message_details)
        self._message_details = prepare_md(
            ts=ts, scores=scores, change_percentage=change_percentage)
        self.send = True

        if self.aid_service_connection is not None:
            message_type = self._message_types[mtype][change_direction]
            logger.debug(
                (repr(self) + ' in _send_message_aid_service_connection:'
                 ' validity:{}, dt:{}, change_to_report:{}, type:{}'.format(
                     self.aid_service_connection, self._datetime_utc,
                     change_percentage, message_type)))
            self._grpc_messages[message_type](*self._message_details)

        else:
            logger.debug(
                (repr(self) + 'in _send_message_aid_service_connection, '
                 'NO aid_service_connection dt: {}, type: {}'.format(
                     self._datetime_utc, mtype)))

    def _check_for_immediate_report(self, model_decision, current_time):
        """
        Check for immediate report.

        Arguments:
        - model_decision: int, the decision of model, 1, 0 or -1 for upward,
                          no trend or downward
        - current_time: Pandas DateTime timestamp

        Returns:
        - trend_decision: is an integer which 1 for growth, -1 for decline
                          and 0 for no trend
        - parameter change_to_report: is float and the amount of change
        """
        logger.debug(
            repr(self) + 'inside _check_for_immediate_report.'
            'model_decision: {} , time:{}'.format(model_decision, current_time)
        )
        empty_verdict = 0, 0
        has_verdict = model_decision != 0
        mismatch_present = (
            not has_verdict
            and self._current_trend_direction + model_decision == 0)

        if not has_verdict:
            logger.debug('No verdict!')
            self._current_trend_direction = 0
            return empty_verdict

        elif mismatch_present:
            logger.error("Trend direction mismatch between DM ({}) and "
                         "new value ({})".format(self._current_trend_direction,
                                                 model_decision))
            return empty_verdict

        self._trends_list.append(model_decision)
        already_in_trend = self._current_trend_direction != 0
        logger.debug('Already in trend {}'.format(already_in_trend))
        if already_in_trend:
            return empty_verdict
        logger.debug('Not already in trend!')
        self._current_trend_direction = model_decision
        change_percentage = self.get_change_metric(
            new_val=self._avg_traffic_new, old_val=self._avg_traffic_old)
        not_monthly_report_day = current_time.day not in self.report_days
        change_is_valid = self._is_valid_change(self._current_trend_direction,
                                                change_percentage)

        if not (change_is_valid and not_monthly_report_day):
            logger.debug('Invalid change or not report day')
            return 0, change_percentage

        trend_decision = self._current_trend_direction
        self._send_message_aid_service_connection(change_percentage,
                                                  trend_decision, 'immediate')
        trend_dir_str = {
            1: 'upwards',
            -1: 'downwards'
        }.get(self._current_trend_direction)
        logger.info(
            str(self) + ' immediate {} trend was sent. '
            'immediate change: {} datetime: {}'.format(
                trend_dir_str, change_percentage, self._datetime_utc))
        return trend_decision, change_percentage

    def _check_for_monthly_report(self):
        """
        Check for monthly reports

        Returns:
        - monthly_decision: int, is 1 if there has been some growth in
                            the past 15 days to be reported. -1 for declines
                            and zero for nothing
        """
        monthly_trend_direction = self._get_monthly_trend_direction()
        self._trends_list.clear()
        empty_verdict = 0, 0

        if not monthly_trend_direction:
            logger.debug(
                repr(self) + '_check_for_monthly_report.'
                ' None of the conditions are met')
            return empty_verdict

        monthly_change = self.get_change_metric(
            new_val=self._avg_traffic_new, old_val=self._avg_traffic_old)
        change_is_valid = self._is_valid_change(monthly_change,
                                                monthly_trend_direction)
        if not change_is_valid:
            return 0, monthly_change

        self._send_message_aid_service_connection(
            monthly_change, monthly_trend_direction, 'monthly')
        trend_dir_str = {
            1: 'upwards',
            -1: 'downwards'
        }.get(monthly_trend_direction)
        logger.info(
            str(self.app_id) + ' monthly {} trend was sent.'
            'monthly change {}, datetime: {}'.format(
                trend_dir_str, int(monthly_change), self._datetime_utc))

        return monthly_trend_direction, monthly_change

    def _preprocess(self, model_output):
        """
        Preprocess the decision maker input.

        Arguments:
        - model_output: dict, the input of the Class coming from Detector

        Returns:
        - model_decision: int, is an integer which 1 for growth, -1 for decline
                        and 0 for no trend
        - no_traffic_flag: bool, if there have been 3 days with no traffic
        - current_time: time extracted from last date_package
        """
        self._send = False
        self._message_details = None
        current_time = model_output['time_model']
        model_decision = model_output['value_model']
        self._avg_traffic_new = model_output['avg_traffic']
        self._avg_traffic_old = model_output['old_average_traffic']
        no_traffic_flag = model_output['no_traffic_flag']
        self._datetime_utc = current_time

        assert model_decision in (
            -1, 0, 1), 'Invalid decision {}'.format(model_decision)
        return model_decision, no_traffic_flag, current_time

    def _handle_immediate_check(self, no_traffic_flag, model_decision,
                                current_time):
        """
        Get immediate reports.

        Returns trend decision and change to report.
        """
        if no_traffic_flag:
            return 0, 0
        return self._check_for_immediate_report(model_decision, current_time)

    def _handle_monthly_check(self, current_time):
        """
        Get immediate reports.

        Returns trend decision and change to report.
        """
        if current_time.day not in self.report_days:
            return 0, 0
        return self._check_for_monthly_report()

    def __repr__(self):
        """Use repr(self) for DEBUG, ERROR, CRITICAL logging levels."""
        return '{} - Metrics {}: '.format(self.app_id, self.metric)

    def __str__(self):
        """Use str(self) for INFO, WARNING logging levels."""
        return '{}: '.format(self.app_id)

    @property
    def days_to_report(self):
        return self._days_to_report

    @property
    def report_days(self):
        return self._report_days

    @property
    def reportable_change_bounds(self):
        return self._reportable_change_bounds
