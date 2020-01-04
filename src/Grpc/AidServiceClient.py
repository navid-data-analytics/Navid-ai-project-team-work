import sys
sys.path.append('service/gen/protos') # noqa
import ai_decision_service_pb2
import ai_decision_service_pb2_grpc
from src.Grpc.conversions import \
    datetimeToGrpctimestamp, \
    grpctimestampToDatetime, \
    dictToGrpcdata, \
    grpcdataToDict
from datetime import datetime, timezone
from src.Grpc.ConnectionClient import ConnectionClient, DataServiceError
import logging

"""
To (re)generate protos run:

In ai-decision folder:
 docker-compose run ai_decision ./service/protos/generate_protos_python.sh
"""

logger = logging.getLogger('root')

DEFAULT_DT = datetime(1971, 1, 1, tzinfo=timezone.utc)
DEFAULT_APPID = 1


class AidServiceClient(ConnectionClient):
    """
    A gRPC client. Handles the protocols.
    """
    def __init__(self, address):
        super(AidServiceClient, self).__init__(
            address,
            connection_timeout=2,
            max_retries=0)

    # AIDecisionStateServiceStub
    def SaveState(self, keyword, state, dt=None, appID=None):
        """
        Save arbitrary state, no checking is performed.
        input:
            keyword: String, make sure its unique
            state: Dict, can have arbitrary values.
                Method caller responsible for usage of later retrieved state
            dt: Datetime, None if unused
            appID: int, None if unused
        returns:
            Exception, None if no error
        """
        if dt is None:
            dt = DEFAULT_DT
        if appID is None:
            appID = DEFAULT_APPID
        try:
            request = ai_decision_service_pb2.StateSaveRequest(
                app_id=appID,
                keyword=keyword,
                data=dictToGrpcdata(state),
                generation_time=datetimeToGrpctimestamp(dt),
            )
        except (TypeError) as e:
            err = DataServiceError('SaveStateRequest', e)
            logger.error(err)
            return None, err

        service = self.getService(
            ai_decision_service_pb2_grpc.AIDecisionStateServiceStub)
        res, e = self.send(
            service.Save,
            request,
            'SaveState',
            reliable=True)
        return e

    def GetState(self, keyword, dt=None, appID=None):
        """
        Get arbitrary state, saved with SaveState.
        input:
            keyword: String, make sure its unique
            dt: Datetime, None if unused
            appID: int, None if unused
        returns:
            Dict (defined by caller itself in SaveState), None if error
        """
        if dt is None:
            dt = DEFAULT_DT
        if appID is None:
            appID = DEFAULT_APPID
        try:
            request = ai_decision_service_pb2.StateGetRequest(
                app_id=appID,
                keyword=keyword,
                generation_time=datetimeToGrpctimestamp(dt),
            )
        except (TypeError) as e:
            err = DataServiceError('StateGetRequest', e)
            logger.error(err)
            return None, err

        service = self.getService(
            ai_decision_service_pb2_grpc.AIDecisionStateServiceStub)
        res, e = self.send(
            service.Get,
            request,
            'GetState',
            reliable=True)
        if e is not None or not res:
            return None
        return grpcdataToDict(res.data)

    # AIDecisionMessageServiceStub
    def _CreateMessage(self, dt, appID, type, version, data):
        """
        Tell AID-E to create a new message in the database.
        input:
            dt: Datetime, the date/time for which the message is produced
            appID: int
            type: String, uniquely identifying the message
            version: int, version of the message
            data: Dict, data corresponding to type and version of the message.
                Entries are defined for each message separately.
        returns:
            Exception, None if no error
        """
        try:
            request = ai_decision_service_pb2.MessageCreateRequest(
                app_id=appID,
                type=type,
                version=version,
                data=dictToGrpcdata(data),
                generation_time=datetimeToGrpctimestamp(dt),
            )
        except (TypeError) as e:
            info = 'MessageCreateRequest ({} v{})'.format(type, version)
            err = DataServiceError(info, e)
            logger.error(err)
            return None, err

        logger.info("gRPC message {} v{} ({}) appID={}: send".format(
            type, version, dt, appID))
        service = self.getService(
            ai_decision_service_pb2_grpc.AIDecisionMessageServiceStub)
        res, e = self.send(
            service.Create,
            request,
            'CreateMessage',
            reliable=True)
        return e

    def ListMessages(self, appID,
                     type="",
                     minVersion=0, maxVersion=0,
                     start=None, end=None):
        """
        Get a stream of messages.
        input:
            appID: int, the appID to retrieve messages for, mandatory
            type: String, type of message
            minVersion: int, minimum version of messages
            maxVersion: int, maximum version of messages
            start: Datetime, the start of the time frame to query, can be None
            end: Datetime, the end of the time frame to query, can be None
        returns:
            generator, if error occured it is handled only after this generator
                is accessed
            contains dict:
                'message', String, the full message
                'appID', int
                'type', String, type of message template used
                'version', int, version of the message template
                'data': dict, data for the message, entries message specific
                'dt': Datetime object, time of generation of message
        """
        try:
            request = ai_decision_service_pb2.MessageListRequest(
                app_id=appID,
                type=type,
                min_version=minVersion,
                max_version=maxVersion,
                generation_time_from=datetimeToGrpctimestamp(start),
                generation_time_to=datetimeToGrpctimestamp(end)
            )
        except (TypeError) as e:
            err = DataServiceError('MessageListRequest', e)
            logger.error(err)
            return None, err

        service = self.getService(
            ai_decision_service_pb2_grpc.AIDecisionMessageServiceStub)
        res, e = self.send(
            service.List,
            request,
            'ListMessages',
            reliable=True)
        if e is not None:
            return None

        # try-except below is needed, otherwise grpc errors are not catched
        # grpc errors are only raised when the generator is accessed
        # by the CALLER of this function
        try:
            for rawEntry in res:
                entry = {
                    'message': rawEntry.message,
                    'appID': rawEntry.app_id,
                    'type': rawEntry.type,
                    'version': rawEntry.version,
                    'data': grpcdataToDict(rawEntry.data),
                    'dt': grpctimestampToDatetime(rawEntry.generation_time),
                }
                yield entry
        except Exception as e:
            logger.error(e)
            return None
