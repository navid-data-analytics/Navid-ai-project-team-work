"""General class for decision making component."""
from src.components import Transceiver
import logging

logger = logging.getLogger('root')


class DecisionMaker(Transceiver):
    """
    A DecisionMaker object decides if input signal should be forwarded.

    The component is a transceiver that takes signal on the input, and decides
    what should be done with it.
    """

    def __init__(self,
                 app_id=None,
                 metric=None,
                 aid_service_connection=None,
                 process=lambda: None):
        """
        Construct a new DecisionMaker instance.

        Arguments:
        - app_id: integer
        - metric: string specifying the metric
        - aid_service_connection: string
        - process: function started upon arrival of data
        """
        self._app_id = app_id
        self._metric = metric
        self._aid_service_connection = aid_service_connection
        logger.debug(
            str(self) + ' Creating DecisionMaker {} object '.format(app_id))
        Transceiver.__init__(self, process=process)
        logger.debug(
            str(self) + ' DecisionMaker {} object created'.format(app_id))

    def _run(self, input_signal):
        raise NotImplementedError('Not Implemented!')

    def _preprocess(self):
        raise NotImplementedError('Not Implemented!')

    def _decide(self, signal):
        raise NotImplementedError('Not Implemented!')

    def _postprocess(self):
        raise NotImplementedError('Not Implemented!')

    @property
    def app_id(self):
        """The AppID of the DecisionMaker object."""
        return self._app_id

    @property
    def metric(self):
        """The metric object monitors."""
        return self._metric

    @property
    def aid_service_connection(self):
        """Connection to AID-E."""
        return self._aid_service_connection
