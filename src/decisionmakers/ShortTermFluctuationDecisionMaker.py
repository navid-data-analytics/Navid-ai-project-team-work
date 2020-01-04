from src.decisionmakers.FluctuationDecisionMaker import FluctuationDecisionMaker #noqa


class ShortTermFluctuationDecisionMaker(FluctuationDecisionMaker):
    def __init__(self,
                 app_id=None,
                 metric='number of calls',
                 aid_service_connection=None,
                 message_lag=11):
        """
        Create decisionmaker object initialized at stable state.

        Arguments:
        - app_id: integer
        - metric: string specifying the metric
        - aid_service_connection: string
        """
        super(ShortTermFluctuationDecisionMaker,
              self).__init__(app_id, metric, aid_service_connection,
                             message_lag)

    def _init_grpc_messages(self):
        """Create message dictionary."""
        if self.aid_service_connection:
            self.grpc_messages = self.aid_service_connection.messages[
                self.metric]['fluctuation_short_term']
        else:
            self.grpc_messages = None
