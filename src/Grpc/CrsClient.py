import logging
import sys
import src.utils.constants as constants
from src.Grpc.ConnectionClient import ConnectionClient, DataServiceError
from src.Grpc.conversions import datetimeToGrpctimestamp,\
                                 grpctimestampToDatetime

sys.path.append('protos')  # noqa
import conference_reports_service_pb2
import conference_reports_service_pb2_grpc

logger = logging.getLogger('root')


class CrsClient(ConnectionClient):
    """
    A gRPC client. Handles the protocols.
    """

    def __init__(self, address):
        super(CrsClient, self).__init__(
            address, connection_timeout=2, max_retries=0)

    def Aggregate(self, appID, from_dt, to_dt):
        """
        Fetch aggregated data
        input:
            appID: int, the appID for which to aggregate
            from_dt: Datetime, the start of the aggregation
            to_dt: Datetime, the end of the aggregation
        returns:
            None if Error
            dict:
                appID: int, the appID for which it was queried
                from_dt: Datetime, the start of aggregated data
                to_dt: Datetime, the end of aggregated data
                conferences_terminated: int, number of conferences
                objective_quality_v35_average: float, avg OQ
                rtt_average: float, avg RTT in ms
        """
        try:
            request = conference_reports_service_pb2.AggregateRequest(
                app_id=int(appID),
                from_ts=datetimeToGrpctimestamp(from_dt),
                to_ts=datetimeToGrpctimestamp(to_dt),
                # conf_id=None,
                # user_id=None,
                # aggregation_filter=None,
                # request_id='',
                # async=False,
            )
        except (TypeError) as e:
            err = DataServiceError('AggregateRequest', e)
            logger.error(err)
            return None, err

        service = self.getService(
            conference_reports_service_pb2_grpc.ConferenceReportsServiceStub)
        res, e = self.send(
            service.Aggregate, request, 'Aggregate', reliable=True)
        if e is not None or not res:
            return None
        res_dict = {
            'appID':
                int(res.app_id),
            'from_dt':
                grpctimestampToDatetime(res.from_ts),
            'to_dt':
                grpctimestampToDatetime(res.to_ts),
            'conferences_terminated':
                int(res.conferences_terminated),
            'objective_quality_v35_average':
                float(res.objective_quality_v35_average),
            'rtt_average':
                float(res.rtt_average.nanos) / constants.CONVERT_GIGA *
                constants.CONVERT_KILO,
            'delayEffectMean': float(res.delay_effect_mean),
            'throughputEffectMean': float(res.throughput_effect_mean),
            'lossEffectMean': float(res.loss_effect_mean),
        }
        return res_dict
