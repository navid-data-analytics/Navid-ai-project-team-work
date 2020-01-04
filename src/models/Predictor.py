"""General class for forecasting models."""
from src.models import Model
from collections import deque
import logging

logger = logging.getLogger('root')


class Predictor(Model):
    """
    A Predictor object uses history to try to forecast the future.

    The component is a transceiver that takes dataframe on the input,
    and once a certain history threshold is filled it starts to predict
    future values. The threshold is set in model_configuration dictionary
    with 'history_size' integer value.
    """

    def __init__(self, model_configuration, app_id=None, metric=None):
        """
        Construct a new Predictor instance.

        Arguments:
        - model_configuration: dictionary containing model_configuration
        - app_id: integer
        """
        super(Predictor, self).__init__(app_id=app_id, process=self._run,
                                        metric=metric)
        logger.debug(str(self) + ' Creating Predictor object')
        logger.debug(str(self) + ' Setting initial configuration of the model')
        self.model_configuration = model_configuration
        self.history = deque(maxlen=model_configuration['history_size'])
        self.get_last_element = lambda x: x[-1]
        logger.debug(str(self) + ' Predictor created')

    def _forecast(self):
        raise NotImplementedError('Not Implemented!')
