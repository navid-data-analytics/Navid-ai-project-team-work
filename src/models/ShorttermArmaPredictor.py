from src.models import ArmaPredictor
from prometheus_client import Gauge
import logging

output_prometheus = Gauge('Shortterm_ArmaPredictor_output', 'appIDs',
                          ['app_id', 'metric_name'])
logger = logging.getLogger('root')
PREDICT_INTERVAL = 7


class ShorttermArmaPredictor(ArmaPredictor):
    def __init__(self,
                 model_configuration,
                 app_id=None,
                 predict_interval=PREDICT_INTERVAL,
                 metric=None):
        """
        Construct a new ArmaPredictor instance for Shortterm.
        """
        super(ShorttermArmaPredictor, self).__init__(
            model_configuration, app_id, predict_interval, metric)
