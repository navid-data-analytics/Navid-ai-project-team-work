import logging
from src.pipeline import AppPipeline
from src.decisionmakers import DecisionMakerMidTerm
from src.models import TrendDetector

logger = logging.getLogger('root')

# aliases
TrendDetectionModel = TrendDetector
TrendDetectionDM = DecisionMakerMidTerm


class AppPipelineEffect(AppPipeline):
    def __init__(self, config, appID, AidServiceConnection, metric):
        super(AppPipelineEffect, self).__init__(config, appID,
                                                AidServiceConnection, metric)

    def _create_components(self, config):
        self._trend_detection_model = TrendDetectionModel(
            config['detection'], app_id=self.appID, metric=self.metric)
        self._trend_detection_decisionmaker = DecisionMakerMidTerm(
            aid_service_connection=self._AidServiceConnection,
            app_id=self.appID,
            metric=self.metric)

    def _joiner_process(self, values):
        return {
            'trend_detection': values['trend_detection'],
        }

    def _link_components(self):
        # connect trend detection
        self._entrypoint.output.connect(self._trend_detection_model.input)
        self._trend_detection_model.output.connect(
            self._trend_detection_decisionmaker.input)
        self._trend_detection_decisionmaker.output.connect(
            self._exitpoint.get_input('trend_detection'))
