import logging
from src.pipeline import AppPipeline
from src.decisionmakers import ShortTermFluctuationDecisionMaker
from src.models import MidTermFluctuationDetector

logger = logging.getLogger('root')

# aliases

FluctuationDetectionModel = MidTermFluctuationDetector
FluctuationDM = ShortTermFluctuationDecisionMaker

METRIC_NAME = 'objective_quality_v35_average'
TREND_AID_SERVICE_CONNECTION = None  # trend detection won't send grpc msg


class AppPipelineOQShortTerm(AppPipeline):
    def __init__(self, config, appID, AidServiceConnection):
        super(AppPipelineOQShortTerm,
              self).__init__(config, appID, AidServiceConnection, METRIC_NAME)

    def _create_components(self, config):
        self._fluctuation_detection = FluctuationDetectionModel(
            config['fluctuation'],
            app_id=self.appID,
            metric=self.metric,
            window_size=7)
        self._fluctuation_decision = FluctuationDM(
            aid_service_connection=self._AidServiceConnection,
            app_id=self.appID,
            metric=self.metric)

    def _joiner_process(self, values):
        return {
            'fluctuation_detection': values['fluctuation_detection'],
            'trend_detection': {
                None: {
                    METRIC_NAME: None
                }
            },
        }

    def _link_components(self):
        # connect fluctuation detection
        self._entrypoint.output.connect(self._fluctuation_detection.input)
        self._fluctuation_detection.output.connect(
            self._fluctuation_decision.input)
        self._fluctuation_decision.output.connect(
            self._exitpoint.get_input('fluctuation_detection'))
