import logging
from src.pipeline import AppPipeline
from src.decisionmakers import DecisionMakerShortTermRtt
from src.decisionmakers import ShortTermFluctuationDecisionMaker
from src.models import MidTermFluctuationDetector, ShortTrendDetector


logger = logging.getLogger('root')

# aliases
TrendDetectionModel = ShortTrendDetector
FluctuationDetectionModel = MidTermFluctuationDetector

TrendDetectionDM = DecisionMakerShortTermRtt
FluctuationDM = ShortTermFluctuationDecisionMaker

METRIC_NAME = 'rtt_average'


class AppPipelineRTTShortTerm(AppPipeline):

    def __init__(self, config, appID, AidServiceConnection):
        super(AppPipelineRTTShortTerm, self).__init__(config, appID,
                                                      AidServiceConnection,
                                                      METRIC_NAME)

    def _create_components(self, config):
        self._trend_detection_model = TrendDetectionModel(config['detection'], #noqa
                                                          app_id=self.appID,
                                                          metric=self.metric)
        self._trend_detection_decisionmaker = TrendDetectionDM(
                aid_service_connection=self._AidServiceConnection,
                app_id=self.appID,
                metric=self.metric)

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
        return {'fluctuation_detection': values['fluctuation_detection'],
                'trend_detection': values['trend_detection']
                }

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
