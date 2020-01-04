from src.decisionmakers import DecisionMakerShortTerm
import logging

logger = logging.getLogger('root')


class DecisionMakerShortTermRtt(DecisionMakerShortTerm):
    """
    This class is responsible for deciding what kind
    and when a notification should be created.
    """

    def __init__(self,
                 aid_service_connection=None,
                 app_id=None,
                 metric=None,
                 report_days=(None, ),
                 reportable_change_bounds=(20, 10000),
                 days_to_report=14):
        """
        Construct a new DecisionMakerShortTerm  for RTT feature.
        It inherits from DecisionMakerShortTerm Class while
        DecisionMakerShortTerm has inherited from the base Class
        DecisionMaker. This class has two additional functions:
        _check_rtt_growth_by_threshold and
        _check_rtt_decline_by_threshold
        and it overrides the methods _check_for_immediate_report
        and _check_for_monthly_report

        Arguments:
        - aid_service_connection: MessageClient object, if None the component
                                  will not send grpc messages
        - app_id: integer
        - metric: string specifying the metric
        - report_days: tuple of ints, which days to send monthly messages at
        - reportable_change_bounds: tuple of ints, what are the absolute
                                    boundaries of a valid change value
        - days_to_report: how many days are included within a report
        """
        super(DecisionMakerShortTermRtt, self).__init__(
            app_id=app_id,
            metric=metric,
            aid_service_connection=aid_service_connection,
            report_days=report_days,
            reportable_change_bounds=reportable_change_bounds,
            days_to_report=days_to_report)

    def get_change_metric(self, new_val, old_val):
        """
        Calculates the change value

        Arguments:
        - new_val: float, new value for measured metric
        - old_val: float, old value for measured metric

        Returns:
         - change: float, difference between new and old value
        """
        return new_val - old_val

    def _is_valid_change(self, rtt_change, direction):
        """
        Check if the change direction and magnitude are valid.

        Arguments:
        - change: int, change value to be validated
        - direction: 1 or -1, for positive and negative slope

        Returns:
        - True if valid, False otherwise
        """
        return (self.reportable_change_bounds[0] < rtt_change * direction <
                self.reportable_change_bounds[1])

    def _prepare_monthly_message_details(self, ts, scores, *args, **kwargs):
        return [self._datetime_utc, ts, scores, self.app_id]

    def _prepare_immediate_message_details(self, ts, scores, *args, **kwargs):
        return [self._datetime_utc, ts, scores, self.app_id]
