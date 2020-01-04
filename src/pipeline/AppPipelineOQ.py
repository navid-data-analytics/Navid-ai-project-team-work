import logging
from src.pipeline import AppPipeline
from src.decisionmakers import FluctuationDecisionMaker, DecisionMakerMidTerm
from src.models import MidTermFluctuationDetector, TrendDetector


logger = logging.getLogger('root')

# aliases
TrendDetectionModel = TrendDetector
TrendDetectionDM = DecisionMakerMidTerm

FluctuationDetectionModel = MidTermFluctuationDetector
FluctuationDM = FluctuationDecisionMaker

METRIC_NAME = 'objective_quality_v35_average'
TREND_AID_SERVICE_CONNECTION = None  # trend detection won't send grpc msg


class AppPipelineOQ(AppPipeline):

    def __init__(self, config, appID, AidServiceConnection):
        super(AppPipelineOQ, self).__init__(config, appID,
                                            AidServiceConnection,
                                            METRIC_NAME)

    def _create_components(self, config):
        self._trend_detection_model = TrendDetectionModel(
            config['detection'], app_id=self.appID, metric=self.metric)
        self._trend_detection_decisionmaker = TrendDetectionDM(
                aid_service_connection=TREND_AID_SERVICE_CONNECTION,
                app_id=self.appID,
                metric=self.metric)

        self._fluctuation_detection = FluctuationDetectionModel(
                                        config['fluctuation'],
                                        app_id=self.appID,
                                        metric=self.metric)
        self._fluctuation_decision = FluctuationDecisionMaker(
                            aid_service_connection=self._AidServiceConnection,
                            app_id=self.appID,
                            metric=self.metric)

    def _joiner_process(self, values):
        return {'fluctuation_detection': values['fluctuation_detection'],
                'trend_detection': values['trend_detection'], }

    def _link_components(self):
        # connect fluctuation detection
        self._entrypoint.output.connect(self._fluctuation_detection.input)
        self._fluctuation_detection.output.connect(
                                    self._fluctuation_decision.input)
        self._fluctuation_decision.output.connect(
            self._exitpoint.get_input('fluctuation_detection'))

        # connect trend detection
        self._entrypoint.output.connect(self._trend_detection_model.input)
        self._trend_detection_model.output.connect(
            self._trend_detection_decisionmaker.input)
        self._trend_detection_decisionmaker.output.connect(
            self._exitpoint.get_input('trend_detection'))
