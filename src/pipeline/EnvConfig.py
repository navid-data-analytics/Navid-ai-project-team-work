import os
import sys
import logging
from src.pipeline import AppPipelineVolume, AppPipelineOQ, AppPipelineRTT
from src.pipeline import AppPipelineVolumeShortTerm, AppPipelineOQShortTerm, AppPipelineRTTShortTerm #noqa
from src.pipeline import AppPipelineDelayEffect, AppPipelineLossEffect
from src.pipeline import AppPipelineThroughputEffect
from src.utils import setup_logger

logger = setup_logger(logging.getLogger('root'),
                      handler=logging.StreamHandler(),
                      level=os.getenv('LOG_LEVEL'))


class EnvConfig:
    def __init__(self):
        self._env = os.getenv('ENV')
        self._sentry_credentials = \
            None if 'http' not in os.getenv('SENTRY_CREDENTIALS') \
            else os.getenv('SENTRY_CREDENTIALS')
        self._version = os.getenv('VERSION')
        self._log_level = os.getenv('LOG_LEVEL')
        self._parse_app_ids()
        self._port_prometheus = int(os.getenv('PORT_PROMETHEUS'))
        self._aid_service_grpc_address = os.getenv('AID_SERVICE_GRPC_ADDR')
        self._crs_grpc_address = os.getenv('CRS_GRPC_ADDR')
        self._supported_terms = ['midterm', 'shortterm']
        self._supported_metric_names = [
                                        # 'delayEffectMean',
                                        # 'throughputEffectMean',
                                        # 'lossEffectMean',
                                        'conferences_terminated',
                                        # 'objective_quality_v35_average',
                                        'rtt_average']
        self._supported_metrics = {'midterm': {
            # 'delayEffectMean': AppPipelineDelayEffect,
            # 'throughputEffectMean': AppPipelineThroughputEffect,
            # 'lossEffectMean': AppPipelineLossEffect,
            'conferences_terminated': AppPipelineVolume,
            # 'objective_quality_v35_average': AppPipelineOQ,
            'rtt_average': AppPipelineRTT},
                                'shortterm': {
            'conferences_terminated': AppPipelineVolumeShortTerm,
            # 'objective_quality_v35_average': AppPipelineOQShortTerm,
            'rtt_average': AppPipelineRTTShortTerm,
            # 'delayEffectMean': AppPipelineDelayEffect,
            # 'throughputEffectMean': AppPipelineThroughputEffect,
            # 'lossEffectMean': AppPipelineLossEffect
                                            },
            }

    def _parse_app_ids(self):
        logger.debug('Parsing appIDs')
        if not os.getenv('APPIDS'):
            logger.warning('Please set the appIDs in docker-compose.yml')
            sys.exit()
        self._app_ids = [
            int(appID) for appID in os.getenv('APPIDS').split(',')
        ]
        logger.debug('appIDs parsed correctly')

    @property
    def env(self):
        return self._env

    @property
    def sentry_credentials(self):
        return self._sentry_credentials

    @property
    def version(self):
        return self._version

    @property
    def app_ids(self):
        return self._app_ids

    @property
    def log_level(self):
        return self._log_level

    @property
    def port_prometheus(self):
        return self._port_prometheus

    @property
    def aid_service_grpc_address(self):
        return self._aid_service_grpc_address

    @property
    def crs_grpc_address(self):
        return self._crs_grpc_address

    @property
    def supported_metrics(self):
        return self._supported_metrics

    @property
    def supported_terms(self):
        return self._supported_terms

    @property
    def supported_metric_names(self):
        return self._supported_metric_names
