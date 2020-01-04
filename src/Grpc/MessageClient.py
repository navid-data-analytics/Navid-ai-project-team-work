from src.Grpc.StateClient import StateClient
from src.utils import to_unix_timestamp


class MessageClient(StateClient):
    def __init__(self,
                 address,
                 flags={
                     'unsuppress': {},
                     'date': {}
                 },
                 load_state=False):
        super(MessageClient, self).__init__(address, flags, load_state)
        self._initialize_message_dict()

    def _initialize_message_dict(self):
        self.messages = {
            'conferences_terminated': {
                'trend': {
                    'Immediate_down': self.CreateVolumeMidtermTrendImmediatelyDown,
                    '15_days_down': self.CreateVolumeMidtermTrend15daysDown,
                    '15_days_up': self.CreateVolumeMidtermTrend15daysUp,
                    'Immediate_up': self.CreateVolumeMidtermTrendImmediatelyUp,
                    'Immediate_down_short': self.CreateVolumeShorttermTrendImmediatelyDown,
                    'Immediate_up_short': self.CreateVolumeShorttermTrendImmediatelyUp,
                         },
                'prediction': {
                    '7_days_down': self.CreateVolumeShorttermPrediction7daysDown,
                    '7_days_up': self.CreateVolumeShorttermPrediction7daysUp,
                    '15_days_down': self.CreateVolumeMidtermPrediction15daysDown,
                    '15_days_up': self.CreateVolumeMidtermPrediction15daysUp
                                },
                'fluctuation': {
                   -1: self.CreateVolumeMidtermFluctuationImmediatelyLow,
                    0: self.CreateVolumeMidtermFluctuationImmediatelyStabilized,
                    1: self.CreateVolumeMidtermFluctuationImmediatelyHigh,
                                },
                'fluctuation_short_term': {
                   -1: self.CreateVolumeShorttermFluctuationImmediatelyLow,
                    0: self.CreateVolumeShorttermFluctuationImmediatelyStabilized,
                    1: self.CreateVolumeShorttermFluctuationImmediatelyHigh,
                                }
                                        },
            'objective_quality_v35_average': {
                # pylama:ignore=E501
                'trend': {
                     1: {  # TREND UP
                        'immediate': {
                            ('lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrendImmediateUpLoss,
                            ('objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrendImmediateUpThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrendImmediateUpDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrendImmediateUpLossDelay,
                            ('lossEffectMean',  'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrendImmediateUpLossThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrendImmediateUpThroughputDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrendImmediateUpAll,
                            },
                        'immediate_short': {
                            ('lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNShorttermTrendImmediateUpLoss,
                            ('objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNShorttermTrendImmediateUpThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average'): self.CreateOQCNShorttermTrendImmediateUpDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNShorttermTrendImmediateUpLossDelay,
                            ('lossEffectMean',  'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNShorttermTrendImmediateUpLossThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNShorttermTrendImmediateUpThroughputDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNShorttermTrendImmediateUpAll,
                            },
                        'monthly': {
                            ('lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrend15daysUpLoss,
                            ('objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrend15daysUpThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrend15daysUpDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrend15daysUpLossDelay,
                            ('lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrend15daysUpLossThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrend15daysUpThroughputDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrend15daysUpAll,
                            },
                        },
                    -1: {  # TREND DOWN
                        'immediate': {
                            ('lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrendImmediateDownLoss,
                            ('objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrendImmediateDownThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrendImmediateDownDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrendImmediateDownLossDelay,
                            ('lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrendImmediateDownLossThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrendImmediateDownThroughputDelay,
                            ('delayEffectMean', 'lossEffectMean',  'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrendImmediateDownAll,
                            },
                        'immediate_short': {
                            ('lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNShorttermTrendImmediateDownLoss,
                            ('objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNShorttermTrendImmediateDownThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average'): self.CreateOQCNShorttermTrendImmediateDownDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNShorttermTrendImmediateDownLossDelay,
                            ('lossEffectMean',  'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNShorttermTrendImmediateDownLossThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNShorttermTrendImmediateDownThroughputDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNShorttermTrendImmediateDownAll,
                            },
                        'monthly': {
                            ('lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrend15daysDownLoss,
                            ('objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrend15daysDownThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrend15daysDownDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average'): self.CreateOQCNMidtermTrend15daysDownLossDelay,
                            ('lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrend15daysDownLossThroughput,
                            ('delayEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrend15daysDownThroughputDelay,
                            ('delayEffectMean', 'lossEffectMean', 'objective_quality_v35_average', 'throughputEffectMean'): self.CreateOQCNMidtermTrend15daysDownAll,
                            },
                        },
                    },
                'fluctuation': {
                   -1: self.CreateOQMidtermFluctuationImmediatelyLow,
                    0: self.CreateOQMidtermFluctuationImmediatelyStabilized,
                    1: self.CreateOQMidtermFluctuationImmediatelyHigh,
                                },
                'fluctuation_short_term': {
                   -1: self.CreateOQShorttermFluctuationImmediatelyLow,
                    0: self.CreateOQShorttermFluctuationImmediatelyStabilized,
                    1: self.CreateOQShorttermFluctuationImmediatelyHigh,
                                },
                    },
            'rtt_average':  {
                'trend': {
                    'Immediate_down': self.CreateRttMidtermTrendImmediatelyDown, #noqa
                    '15_days_down': self.CreateRttMidtermTrend15daysDown,
                    '15_days_up': self.CreateRttMidtermTrend15daysUp,
                    'Immediate_up': self.CreateRttMidtermTrendImmediatelyUp,
                    'Immediate_down_short': self.CreateRttShorttermTrendImmediatelyDown,
                    'Immediate_up_short': self.CreateRttShorttermTrendImmediatelyUp,
                         },
                'fluctuation': {
                   -1: self.CreateRttMidtermFluctuationImmediatelyLow,
                    0: self.CreateRttMidtermFluctuationImmediatelyStabilized,
                    1: self.CreateRttMidtermFluctuationImmediatelyHigh,
                                },
                'fluctuation_short_term': {
                   -1: self.CreateRttShorttermFluctuationImmediatelyLow,
                    0: self.CreateRttShorttermFluctuationImmediatelyStabilized,
                    1: self.CreateRttShorttermFluctuationImmediatelyHigh,
                                },
                    }

                }

    # TREND
    def CreateVolumeMidtermTrend15daysUp(self, dt, timestamps, scores, appID,
                                         percentage):
        """
        Example message: "Your daily average calls have \
                          increased by 5% in the last 30 days. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: self explainatory
            percentage: int, percentage of the growth

        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermTrend15daysUp"
        version = 2
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
            'percentage':
            int(percentage),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeMidtermTrend15daysDown(self, dt, timestamps, scores, appID,
                                           percentage):
        """
        Example message: "Your daily average calls have
                          decreased by 5% in the last 30 days.
                          This could be worth looking into."
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: self explainatory
            percentage: int, percentage of the growth
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermTrend15daysDown"
        version = 2
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
            'percentage':
            int(percentage),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeMidtermTrendImmediatelyUp(self, dt, timestamps, scores,
                                              appID, percentage, days):
        """
        Example message: "Your daily average calls have been growing
                          at a rate of 8% for the past 12 days. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: int
            percentage: int, percentage of the growth
            days: int, number of days
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermTrendImmediatelyUp"
        version = 2
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
            'percentage':
            int(percentage),
            'days':
            int(days),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeMidtermTrendImmediatelyDown(self, dt, timestamps, scores,
                                                appID, percentage, days):
        """
        Example message: "Your daily average calls is decreasing at a rate
                          of 3% for the past 13 days. This could be worth
                          looking into."
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: int
            percentage: int, percentage of the decline
            days: int, number of days
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermTrendImmediatelyDown"
        version = 2
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
            'percentage':
            int(percentage),
            'days':
            int(days),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeShorttermTrendImmediatelyUp(self, dt, timestamps, scores,
                                                appID, percentage, days):
        """
        Example message: "Your daily average calls have been growing
                          at a rate of 8% for the past 7 days. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: int
            percentage: int, percentage of the growth
            days: int, number of days
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermTrendImmediatelyUp"
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
            'percentage':
            int(percentage),
            'days':
            int(days),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeShorttermTrendImmediatelyDown(self, dt, timestamps, scores,
                                                  appID, percentage, days):
        """
        Example message: "Your daily average calls is decreasing at a rate
                          of 3% for the past 7 days. This could be worth
                          looking into."
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: int
            percentage: int, percentage of the decline
            days: int, number of days
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermTrendImmediatelyDown"
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
            'percentage':
            int(percentage),
            'days':
            int(days),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateRttShorttermTrendImmediatelyDown(self, dt, timestamps, scores,
                                               appID):
        """
        Example message:

        "Your avg. Round Trip Time (RTT) per day has decreased.

        Previous period (6th Apr - 6th May) avg. RTT/day was 200 ms.
        Current period (6th May - 6th Jun) avg. RTT/day is 180 ms.

        Great job!"

        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: self explainatory
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermRttTrendImmediatelyDown"
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateRttShorttermTrendImmediatelyUp(self, dt, timestamps, scores,
                                             appID):
        """
        Example message:

        "Your avg. Round Trip Time (RTT) per day has decreased.

        Previous period (6th Apr - 6th May) avg. RTT/day was 200 ms.
        Current period (6th May - 6th Jun) avg. RTT/day is 180 ms.

        Great job!"

        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: self explainatory
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermRttTrendImmediatelyUp"
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    # NOTE: OQ Complex Messages
    def CreateOQCNMidtermTrend15daysUpLoss(self, dt, timestamps, scores, appID,
                                           change_to_report):
        type = 'MidtermOQCNTrend15daysUpLoss'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysUpThroughput(self, dt, timestamps, scores,
                                                 appID, change_to_report):
        type = 'MidtermOQCNTrend15daysUpThroughput'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysUpDelay(self, dt, timestamps, scores,
                                            appID, change_to_report):
        type = 'MidtermOQCNTrend15daysUpDelay'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysUpLossDelay(self, dt, timestamps, scores,
                                                appID, change_to_report):
        type = 'MidtermOQCNTrend15daysUpLossDelay'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysUpLossThroughput(
            self, dt, timestamps, scores, appID, change_to_report):
        type = 'MidtermOQCNTrend15daysUpLossThroughput'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysUpThroughputDelay(
            self, dt, timestamps, scores, appID, change_to_report):
        type = 'MidtermOQCNTrend15daysUpThroughputDelay'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysUpAll(self, dt, timestamps, scores, appID,
                                          change_to_report):
        type = 'MidtermOQCNTrend15daysUpAll'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysDownLoss(self, dt, timestamps, scores,
                                             appID, change_to_report):
        type = 'MidtermOQCNTrend15daysDownLoss'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysDownDelay(self, dt, timestamps, scores,
                                              appID, change_to_report):
        type = 'MidtermOQCNTrend15daysDownDelay'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysDownThroughput(
            self, dt, timestamps, scores, appID, change_to_report):
        type = 'MidtermOQCNTrend15daysDownThroughput'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysDownLossDelay(self, dt, timestamps, scores,
                                                  appID, change_to_report):
        type = 'MidtermOQCNTrend15daysDownLossDelay'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysDownLossThroughput(
            self, dt, timestamps, scores, appID, change_to_report):
        type = 'MidtermOQCNTrend15daysDownLossThroughput'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysDownThroughputDelay(
            self, dt, timestamps, scores, appID, change_to_report):
        type = 'MidtermOQCNTrend15daysDownThroughputDelay'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrend15daysDownAll(self, dt, timestamps, scores,
                                            appID, change_to_report):
        type = 'MidtermOQCNTrend15daysDownAll'
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(change_to_report),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQCNMidtermTrendImmediateUpLoss(self, dt, timestamps, scores,
                                              appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateUpLoss'
        return type

    def CreateOQCNMidtermTrendImmediateUpThroughput(
            self, dt, timestamps, scores, appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateUpThroughput'
        return type

    def CreateOQCNMidtermTrendImmediateUpDelay(self, dt, timestamps, scores,
                                               appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateUpDelay'
        return type

    def CreateOQCNMidtermTrendImmediateUpLossDelay(
            self, dt, timestamps, scores, appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateUpLossDelay'
        return type

    def CreateOQCNMidtermTrendImmediateUpLossThroughput(
            self, dt, timestamps, scores, appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateUpLossThroughput'
        return type

    def CreateOQCNMidtermTrendImmediateUpThroughputDelay(
            self, dt, timestamps, scores, appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateUpThroughputDelay'
        return type

    def CreateOQCNMidtermTrendImmediateUpAll(self, dt, timestamps, scores,
                                             appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateUpAll'
        return type

    def CreateOQCNMidtermTrendImmediateDownLoss(self, dt, timestamps, scores,
                                                appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateDownLoss'
        return type

    def CreateOQCNMidtermTrendImmediateDownDelay(self, dt, timestamps, scores,
                                                 appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateDownDelay'
        return type

    def CreateOQCNMidtermTrendImmediateDownThroughput(
            self, dt, timestamps, scores, appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateDownThroughput'
        return type

    def CreateOQCNMidtermTrendImmediateDownLossDelay(
            self, dt, timestamps, scores, appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateDownLossDelay'
        return type

    def CreateOQCNMidtermTrendImmediateDownLossThroughput(
            self, dt, timestamps, scores, appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateDownLossThroughput'
        return type

    def CreateOQCNMidtermTrendImmediateDownThroughputDelay(
            self, dt, timestamps, scores, appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateDownThroughputDelay'
        return type

    def CreateOQCNMidtermTrendImmediateDownAll(self, dt, timestamps, scores,
                                               appID, percentage, days):
        """Dummy, to be implemented"""
        type = 'MidtermOQCNTrendImmediateDownAll'
        return type

    def CreateOQCNShorttermTrendImmediateUpLoss(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateUpLoss'
        return type

    def CreateOQCNShorttermTrendImmediateUpThroughput(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateUpThroughput'
        return type

    def CreateOQCNShorttermTrendImmediateUpDelay(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateUpDelay'
        return type

    def CreateOQCNShorttermTrendImmediateUpLossDelay(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateUpLossDelay'
        return type

    def CreateOQCNShorttermTrendImmediateUpLossThroughput(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateUpLossThroughput'
        return type

    def CreateOQCNShorttermTrendImmediateUpThroughputDelay(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateUpThroughputDelay'
        return type

    def CreateOQCNShorttermTrendImmediateUpAll(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateUpAll'
        return type

    def CreateOQCNShorttermTrendImmediateDownLoss(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateDownLoss'
        return type

    def CreateOQCNShorttermTrendImmediateDownThroughput(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateDownThroughput'
        return type

    def CreateOQCNShorttermTrendImmediateDownDelay(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateDownDelay'
        return type

    def CreateOQCNShorttermTrendImmediateDownLossDelay(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateDownLossDelay'
        return type

    def CreateOQCNShorttermTrendImmediateDownLossThroughput(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateDownLossThroughput'
        return type

    def CreateOQCNShorttermTrendImmediateDownThroughputDelay(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateDownThroughputDelay'
        return type

    def CreateOQCNShorttermTrendImmediateDownAll(self):
        """Dummy, to be implemented"""
        type = 'ShorttermOQCNTrendImmediateDownAll'
        return type

    def CreateOQMidtermTrend15daysUp(self, dt, timestamps, scores, appID,
                                     percentage):
        """
        Example message:

        "Your avg. objective quality per day has decreased by 5%.

        Previous period (6th Apr - 6th May) avg. OQ/day was 2.10
        Current period (6th May - 6th Jun) avg. OQ/day is 2.00

        Great job!"

        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: self explainatory
            percentage: int, percentage of the growth

        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermOQTrend15daysUp"
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(percentage),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQMidtermTrend15daysDown(self, dt, timestamps, scores, appID,
                                       percentage):
        """
        Example message:

        "Your avg. objective quality per day has decreased by 5%.

        Previous period (6th Apr - 6th May) avg. OQ/day was 2.10
        Current period (6th May - 6th Jun) avg. OQ/day is 2.00

        This could be worth looking into."

        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: self explainatory
            percentage: int, percentage of the growth
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermOQTrend15daysDown"
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            round(previous_score, 2),
            'current_score':
            round(current_score, 2),
            'percentage':
            int(percentage),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateRttMidtermTrend15daysUp(self, dt, timestamps, scores, appID):
        """
        Example message:

        "Your avg. Round Trip Time (RTT) per day has increased.

        Previous period (6th Apr - 6th May) avg. RTT/day was 180 ms.
        Current period (6th May - 6th Jun) avg. RTT/day is 200 ms.

        This could be worth looking into."

        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: self explainatory
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermRttTrend15daysUp"
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateRttMidtermTrend15daysDown(self, dt, timestamps, scores, appID):
        """
        Example message:

        "Your avg. Round Trip Time (RTT) per day has decreased.

        Previous period (6th Apr - 6th May) avg. RTT/day was 200 ms.
        Current period (6th May - 6th Jun) avg. RTT/day is 180 ms.

        Great job!"

        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: self explainatory
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermRttTrend15daysDown"
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateRttMidtermTrendImmediatelyUp(self, dt, timestamps, scores,
                                           appID):
        """
        Example message:

        "Your avg. Round Trip Time (RTT) per day has decreased.

        Previous period (6th Apr - 6th May) avg. RTT/day was 200 ms.
        Current period (6th May - 6th Jun) avg. RTT/day is 180 ms.

        Great job!"

        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: self explainatory
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermRttTrendImmediatelyUp"
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateRttMidtermTrendImmediatelyDown(self, dt, timestamps, scores,
                                             appID):
        """
        Example message:

        "Your avg. Round Trip Time (RTT) per day has decreased.

        Previous period (6th Apr - 6th May) avg. RTT/day was 200 ms.
        Current period (6th May - 6th Jun) avg. RTT/day is 180 ms.

        Great job!"

        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            scores: tuple with previous and current model score
            appID: self explainatory
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermRttTrendImmediatelyDown"
        version = 1
        previous_score, current_score = scores
        previous_period = timestamps['previous_period']
        current_period = timestamps['current_period']
        data = {
            'previous_period_start':
            int(to_unix_timestamp(previous_period['start'])),
            'previous_period_end':
            int(to_unix_timestamp(previous_period['end'])),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'previous_score':
            int(previous_score),
            'current_score':
            int(current_score),
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    # FLUCTUATION
    def CreateVolumeMidtermFluctuationImmediatelyHigh(self, dt, timestamp,
                                                      appID):
        """
        Example message: "Your service usage has been fluctuating in the
                          past 30 days. This could be worth looking into"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermFluctuationImmediatelyHigh"
        version = 2
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeMidtermFluctuationImmediatelyStabilized(
            self, dt, timestamp, appID):
        """
        Example message: "Your service usage fluctuations have been
                          successfully stabilized. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermFluctuationImmediatelyStabilized"
        version = 2
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start'])),
            'current_period_end':
            int(to_unix_timestamp(timestamp['current_period']['end']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeMidtermFluctuationImmediatelyLow(self, *args, **kwargs):
        """Dummy, to be implemented"""
        return None

    def CreateOQMidtermFluctuationImmediatelyHigh(self, dt, timestamp, appID):
        """
        Example message: "Your avg. OQ/day has been fluctuating in the
                          past 30 days. This could be worth looking into"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermOQFluctuationImmediatelyHigh"
        version = 1
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQMidtermFluctuationImmediatelyStabilized(
            self, dt, timestamp, appID):
        """
        Example message: "Your avg. OQ/day fluctuations have been
                          successfully stabilized. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermOQFluctuationImmediatelyStabilized"
        version = 1
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start'])),
            'current_period_end':
            int(to_unix_timestamp(timestamp['current_period']['end']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQMidtermFluctuationImmediatelyLow(self, *args, **kwargs):
        """Dummy, to be implemented"""
        return None

    def CreateRttMidtermFluctuationImmediatelyLow(self, *args, **kwargs):
        """Dummy, to be implemented"""
        return None

    def CreateRttMidtermFluctuationImmediatelyHigh(self, dt, timestamp, appID):
        """
        Example message: "Your avg. RTT/day has been fluctuating in the
                          past 30 days. This could be worth looking into"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermRttFluctuationImmediatelyHigh"
        version = 1
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateRttMidtermFluctuationImmediatelyStabilized(
            self, dt, timestamp, appID):
        """
        Example message: "Your avg. RTT/day fluctuations have been
                          successfully stabilized. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermRttFluctuationImmediatelyStabilized"
        version = 1
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start'])),
            'current_period_end':
            int(to_unix_timestamp(timestamp['current_period']['end']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQShorttermFluctuationImmediatelyLow(self, *args, **kwargs):
        """Dummy, to be implemented"""
        return None

    def CreateOQShorttermFluctuationImmediatelyStabilized(
            self, dt, timestamp, appID):
        """
        Example message: "Your service usage fluctuations have been
                          successfully stabilized. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermOQFluctuationImmediatelyStabilized"
        version = 1
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start'])),
            'current_period_end':
            int(to_unix_timestamp(timestamp['current_period']['end']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateOQShorttermFluctuationImmediatelyHigh(self, dt, timestamp,
                                                    appID):
        """
        Example message: "Your service usage has been fluctuating in the
                          past 30 days. This could be worth looking into"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermOQFluctuationImmediatelyHigh"
        version = 1
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateRttShorttermFluctuationImmediatelyLow(self, *args, **kwargs):
        """Dummy, to be implemented"""
        return None

    def CreateRttShorttermFluctuationImmediatelyStabilized(
            self, dt, timestamp, appID):
        """
        Example message: "Your service usage fluctuations have been
                          successfully stabilized. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermRttFluctuationImmediatelyStabilized"
        version = 1
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start'])),
            'current_period_end':
            int(to_unix_timestamp(timestamp['current_period']['end']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateRttShorttermFluctuationImmediatelyHigh(self, dt, timestamp,
                                                     appID):
        """
        Example message: "Your service usage has been fluctuating in the
                          past 30 days. This could be worth looking into"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermRttFluctuationImmediatelyHigh"
        version = 1
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeShorttermFluctuationImmediatelyLow(self, *args, **kwargs):
        """Dummy, to be implemented"""
        return None

    def CreateVolumeShorttermFluctuationImmediatelyStabilized(
            self, dt, timestamp, appID):
        """
        Example message: "Your service usage fluctuations have been
                          successfully stabilized. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermFluctuationImmediatelyStabilized"
        version = 1
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start'])),
            'current_period_end':
            int(to_unix_timestamp(timestamp['current_period']['end']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeShorttermFluctuationImmediatelyHigh(
            self, dt, timestamp, appID):
        """
        Example message: "Your service usage has been fluctuating in the
                          past 30 days. This could be worth looking into"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with the beginning of fluctuation
            appID: int
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermFluctuationImmediatelyHigh"
        version = 1
        data = {
            'current_period_start':
            int(to_unix_timestamp(timestamp['current_period']['start']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    # PREDICTIONS
    def CreateVolumeMidtermPrediction15daysUp(self, dt, timestamps, appID,
                                              percentage):
        """
        Example message: "Your daily average calls is expected to grow by 6%
                          in the next 30 days. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            appID: int
            timestamps: dictionary with previous and current timestamps
            percentage: int, percentage of the growth
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermPrediction15daysUp"
        version = 2
        current_period = timestamps['current_period']
        future_period = timestamps['future_period']
        data = {
            'percentage':
            int(percentage),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'future_period_start':
            int(to_unix_timestamp(future_period['start'])),
            'future_period_end':
            int(to_unix_timestamp(future_period['end']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeMidtermPrediction15daysDown(self, dt, timestamps, appID,
                                                percentage):
        """
        Example message: "Your daily average calls is expected to decline
                          by 7% in the next 30 days. You might want to check
                          this out"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            appID: int
            percentage: int, percentage of the decline
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "MidtermPrediction15daysDown"
        version = 2
        current_period = timestamps['current_period']
        future_period = timestamps['future_period']
        data = {
            'percentage':
            int(percentage),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'future_period_start':
            int(to_unix_timestamp(future_period['start'])),
            'future_period_end':
            int(to_unix_timestamp(future_period['end']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeShorttermPrediction7daysDown(self, dt, timestamps, appID,
                                                 percentage):
        """
        Example message: "Your daily average calls is expected to decline
                          by 7% in the next 7 days. You might want to check
                          this out"
        input:
            dt: Datetime, the dt given by scheduler
            timestamps: dictionary with previous and current timestamps
            appID: int
            percentage: int, percentage of the decline
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermPrediction7daysDown"
        version = 1
        current_period = timestamps['current_period']
        future_period = timestamps['future_period']
        data = {
            'percentage':
            int(percentage),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'future_period_start':
            int(to_unix_timestamp(future_period['start'])),
            'future_period_end':
            int(to_unix_timestamp(future_period['end']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)

    def CreateVolumeShorttermPrediction7daysUp(self, dt, timestamps, appID,
                                               percentage):
        """
        Example message: "Your daily average calls is expected to grow by 6%
                          in the next 7 days. Great job!"
        input:
            dt: Datetime, the dt given by scheduler
            appID: int
            timestamps: dictionary with previous and current timestamps
            percentage: int, percentage of the growth
        returns:
            (see AidServiceClient._CreateMessage)
        """
        type = "ShorttermPrediction7daysUp"
        version = 1
        current_period = timestamps['current_period']
        future_period = timestamps['future_period']
        data = {
            'percentage':
            int(percentage),
            'current_period_start':
            int(to_unix_timestamp(current_period['start'])),
            'current_period_end':
            int(to_unix_timestamp(current_period['end'])),
            'future_period_start':
            int(to_unix_timestamp(future_period['start'])),
            'future_period_end':
            int(to_unix_timestamp(future_period['end']))
        }
        return self._CreateMessageConsiderState(dt, appID, type, version, data)
