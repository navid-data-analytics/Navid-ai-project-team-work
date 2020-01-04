from src.pipeline import AppPipelineEffect

METRIC_NAME = 'throughputEffectMean'
AID_SERVICE_CONNECTION = None  # this AppPipeline won't send grpc messages


class AppPipelineThroughputEffect(AppPipelineEffect):
    def __init__(self, config, appID, *args):
        super(AppPipelineThroughputEffect, self).__init__(
                                                    config,
                                                    appID,
                                                    AID_SERVICE_CONNECTION,
                                                    METRIC_NAME)
