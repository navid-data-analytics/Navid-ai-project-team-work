from src.pipeline import AppPipelineEffect

METRIC_NAME = 'delayEffectMean'
AID_SERVICE_CONNECTION = None  # this AppPipeline won't send grpc messages


class AppPipelineDelayEffect(AppPipelineEffect):
    def __init__(self, config, appID, *args):
        super(AppPipelineDelayEffect, self).__init__(config,
                                                     appID,
                                                     AID_SERVICE_CONNECTION,
                                                     METRIC_NAME)
