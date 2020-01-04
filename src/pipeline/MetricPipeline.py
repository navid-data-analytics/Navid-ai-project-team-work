import logging
from src.components import Splitter, Joiner

logger = logging.getLogger('root')


class MetricPipeline:
    def __init__(self,
                 metric_app_pipeline,
                 app_ids,
                 config,
                 AidServiceConnection=None):
        """
        Responsible for processing the data for one metric.

        Depending on the Metric, different class inheriting from this one
        should be created. Only thing changed should be the components
        for processing values changed. Overwrite methods _setup_models
        and _setup_decisionmakers with logic pertaining to the metric.

        The following is replicated for each appID for given metric:

                    /- AppPipelineMetric ------i
                   /                            I
        Splitter ---    AppPipelineMetric        ----- Joiner
                  I                            /
                  I___ AppPipelineMetric_____/
        """
        self._AidServiceConnection = AidServiceConnection
        self._entrypoint = Splitter(process=lambda values: values)
        self._metric = metric_app_pipeline.metric
        self._metric_pipeline = {
            appID: metric_app_pipeline(config, appID, AidServiceConnection)
            for appID in app_ids
        }
        self._exitpoint = Joiner(
            process=lambda values:
            {appID: values[appID] for appID in app_ids}
        )
        self.link_metric_pipeline(app_ids)
        self.input = self._entrypoint.input
        self.output = self._exitpoint.output

    def link_metric_pipeline(self, app_ids):
        for appID in app_ids:
            self._entrypoint.get_output(appID).connect(
                self._metric_pipeline[appID].input)
            self._metric_pipeline[appID].output.connect(
                self._exitpoint.get_input(appID))

    @property
    def metric(self):
        return self._metric

    @property
    def metric_pipeline(self):
        self._metric_pipeline

    def __repr__(self):
        """Use repr(self) for DEBUG, ERROR, CRITICAL logging levels."""
        return 'MetricPipeline {} - BaseClass: '.format(self.metric)

    def __str__(self):
        """Use str(self) for INFO, WARNING logging levels."""
        return 'MetricPipeline {}: '.format(self.metric)
