"""Class for Arma decision making component."""
from src.decisionmakers import ArmaDecisionMaker
from src.utils import measure_time_metric, get_time_dict
from prometheus_client import Gauge
import numpy as np
import logging

output_prometheus = Gauge('ShorttermArmaDecisionMaker_output', 'appIDs',
                          ['app_id', 'metric_name'])
rate_of_change_prometheus = Gauge('Shortterm_Prediction_rate_of_change',
                                  'appIDs', ['app_id', 'metric_name'])
logger = logging.getLogger('root')


class ShortArmaDecisionMaker(ArmaDecisionMaker):
    """
    A ShorttermArmaDecisionMaker object decides what verdict to send.

    The component is a DecisionMaker that takes model output as input and
    decides what is the message sent forward. Additionally it sends predicted
    signal for Grafana.
    """

    def __init__(self,
                 app_id=None,
                 metric='number of calls',
                 aid_service_connection=None,
                 predict_interval=7):
        """
        Construct a new ShorttermArmaDecisionMaker instance.

        Arguments:
        - app_id: integer
        - metric: string specifying the metric
        - aid_service_connection: string
        - predict_interval: int, number of days to predict forward
        """
        super(ShortArmaDecisionMaker, self).__init__(
            app_id, metric, aid_service_connection, predict_interval)

    def boundary_check(self, time):
        return time.day in (7, 14, 21, 28)

    def _send_message(self, time, periods):
        """
        Send message to grpc.

        Decide whether Growth or Decline is sent to grpc, send the message.
        Arguments:
        - time: Time in DateTime format
        - periods: Dictionary of time periods with accordance to documentation
        """
        percentage = self._prepare_percentage()
        if percentage == 0:
            logger.info(
                "Predicting growth, message NOT sent (percentage == 0)")
            return

        if percentage > 0:
            self.grpc_messages['7_days_up'](time, periods, self.app_id,
                                            percentage)
            logger.info(
                str(self) + "Predicting growth, message sent with: " +
                "time: {}\ntime_range{}\nappID: {}\npercentage: {}".format(
                    time, periods, self.app_id, percentage))

        if percentage < 0:
            self.grpc_messages['7_days_down'](time, periods, self.app_id,
                                              percentage)
            logger.info(
                str(self) + "Predicting decline, message sent with: " +
                "time: {}\ntime_range{}\nappID: {}\npercentage: {}".format(
                    time, periods, self.app_id, percentage))
