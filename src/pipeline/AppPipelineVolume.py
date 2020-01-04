import logging
from src.pipeline import AppPipeline
from src.decisionmakers import DecisionMakerMidTerm, FluctuationDecisionMaker
from src.decisionmakers import MidtermArmaDecisionMaker
from src.models import MidtermArmaPredictor, TrendDetector
from src.models import PeriodicityRemovalMidtermFluctuationDetector

logger = logging.getLogger('root')

# aliases
PredictionModel = MidtermArmaPredictor
TrendDetectionModel = TrendDetector
FluctuationDetectionModel = PeriodicityRemovalMidtermFluctuationDetector

PredictionDM = MidtermArmaDecisionMaker
TrendDetectionDM = DecisionMakerMidTerm
FluctuationDM = FluctuationDecisionMaker

METRIC_NAME = 'conferences_terminated'


class AppPipelineVolume(AppPipeline):
    def __init__(self, config, appID, AidServiceConnection):
        super(AppPipelineVolume,
              self).__init__(config, appID, AidServiceConnection, METRIC_NAME)

    def _create_components(self, config):
        self._trend_detection_model = TrendDetectionModel(
            config['detection'], app_id=self.appID, metric=self.metric)
        self._trend_detection_decisionmaker = DecisionMakerMidTerm(
            aid_service_connection=self._AidServiceConnection,
            app_id=self.appID,
            metric=self.metric)

        self._fluctuation_detection = FluctuationDetectionModel(
            config['fluctuation'], app_id=self.appID, metric=self.metric)
        self._fluctuation_decision = FluctuationDecisionMaker(
            aid_service_connection=self._AidServiceConnection,
            app_id=self.appID,
            metric=self.metric)

        self._trend_prediction = PredictionModel(
            config['prediction'], app_id=self.appID, metric=self.metric)
        self._trend_prediction_decisionmaker =\
            PredictionDM(app_id=self.appID,
                         aid_service_connection=self._AidServiceConnection,
                         metric=self.metric)

    def _joiner_process(self, values):
        return {
            'fluctuation_detection': values['fluctuation_detection'],
            'trend_detection': values['trend_detection'],
            'trend_prediction': values['trend_prediction'],
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

        # connect trend prediction
        self._entrypoint.output.connect(self._trend_prediction.input)
        self._trend_prediction.output.connect(
            self._trend_prediction_decisionmaker.input)
        self._trend_prediction_decisionmaker.output.connect(
            self._exitpoint.get_input('trend_prediction'))
