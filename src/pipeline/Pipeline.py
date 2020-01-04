import logging
import src.utils.constants as constants
from src.pipeline import Scheduler, DataFetcher
from src.pipeline import DateQueue
from src.pipeline import Notifier, configs, MetricPipeline
from src.decisionmakers import ComplexOQDecisionMaker
from src.components import Joiner
from itertools import product

logger = logging.getLogger('root')


class Pipeline:
    """
    A Pipeline class is responsible for execution of the pipeline.
    Pipeline object should be constructed with config and env arguments
    holding info for dataprocessor, model and environment.
    """

    def __init__(self, env, AidServiceConnection, CrsConnection):
        self._config = configs.PipelineConfig(env)
        self._create_components(env, AidServiceConnection, CrsConnection)
        self._initialize_components(env)

    def _create_components(self, env, AidServiceConnection, CrsConnection):
        logger.info('Creating components')
        self._scheduler = Scheduler(
            trigger_interval=constants.DAY,
            initial_load=self._config['Scheduler_params']['initial_load'])
        self._datequeue = DateQueue()
        self._datafetcher = DataFetcher(env.supported_metric_names,
                                        env.app_ids, CrsConnection)

        self._metric_pipeline = {
            '{}_{}'.format(metric, term): MetricPipeline(
                env.supported_metrics[term][metric], env.app_ids,
                self._config['models'][term][metric], AidServiceConnection)
            for metric, term in product(env.supported_metric_names, env.
                                        supported_terms)
        }

        self._metric_joiner = Joiner(
            process=lambda values: {'{}_{}'.format(
                metric, term): values['{}_{}'.format(metric, term)]
                for metric, term in product(env.supported_metric_names,
                                            env.supported_terms)}

        )

        self._complex_decisionmaker_oq = ComplexOQDecisionMaker(
            app_ids=env.app_ids, aid_service_connection=AidServiceConnection)

        self._notifier = Notifier(
            env.sentry_credentials,
            release=env.version,
            aid_service_connection=AidServiceConnection)
        logger.info('Components created successfully')

    def _initialize_components(self, env):
        logger.info('Initializing components')
        self._scheduler.output.connect(self._datequeue.input)
        self._datequeue.output.connect(self._datafetcher.input)
        for metric in env.supported_metric_names:
            for term in env.supported_terms:
                self._datafetcher.get_output('{}_{}'.format(
                    metric,
                    term)).connect(self._metric_pipeline['{}_{}'.format(
                        metric, term)].input)

                self._metric_pipeline['{}_{}'.format(
                    metric, term)].output.connect(
                        self._metric_joiner.get_input('{}_{}'.format(
                            metric, term)))

        self._metric_joiner.output.connect(
            self._notifier.input)
        # self._complex_decisionmaker_oq.output.connect(self._notifier.input)
        logger.info('Components initialized successfully')

    def start(self):
        self._scheduler.start()

    def stop(self):
        self._scheduler.shutdown()
        self._datequeue.shutdown()
        logger.info('Closed successfully')
