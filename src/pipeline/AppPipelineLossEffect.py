from src.pipeline import AppPipelineEffect

METRIC_NAME = 'lossEffectMean'
AID_SERVICE_CONNECTION = None  # this AppPipeline won't send grpc messages


class AppPipelineLossEffect(AppPipelineEffect):
    def __init__(self, config, appID, *args):
        super(AppPipelineLossEffect, self).__init__(config,
                                                    appID,
                                                    AID_SERVICE_CONNECTION,
                                                    METRIC_NAME)
