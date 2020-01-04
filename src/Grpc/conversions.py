from google.protobuf.timestamp_pb2 import Timestamp
import json
from datetime import timezone, datetime
import logging

logger = logging.getLogger('root')

DATETIME_FORMAT = '%Y-%m-%dT%H:%M:%S.%f %z'


# Conversion functions
def datetimeToGrpctimestamp(dt):
    """
    Converts a Datetime object into a timestamp transmittable via gRPC.
    Make sure dt is UTC timezone and has the correct granularity.
    """
    if dt is None:
        return None
    if dt.tzinfo != timezone.utc:
        info = '_datetimeToGrpctimestamp: {} is not UTC'.format(dt)
        logger.error(info)
        return None
    # Timestamp.FromDateTime has no tzinfo, so we have to remove it
    dtNoTimezone = dt.replace(tzinfo=None)

    grpctimestamp = Timestamp()
    grpctimestamp.FromDatetime(dtNoTimezone)
    return grpctimestamp


def grpctimestampToDatetime(grpctimestamp):
    """
    Converts a time received over grpc to a Datetime object (UTC).
    """
    dt = grpctimestamp.ToDatetime()  # this is always UTC
    dtTimezone = dt.replace(tzinfo=timezone.utc)
    return dtTimezone


def json_serial(obj):
    """JSON serializer for objects not serializable by default json code"""

    if isinstance(obj, datetime):
        if obj.tzinfo is None:
            raise ValueError("Datetime object does not have tzinfo")
        return obj.strftime(DATETIME_FORMAT)
    raise TypeError("Type %s not serializable" % type(obj))


def dictToGrpcdata(data):
    """
    Converts a Dictionary into bytes data transmittable via gRPC.
    """
    datastr = json.dumps(data, default=json_serial)
    bytes = datastr.encode('utf-8')
    return bytes


def json_deserial(json_dict):
    """JSON de-serializer for objects not de-serializable
        by default json code"""
    for (key, value) in json_dict.items():
        try:
            json_dict[key] = datetime.strptime(value, DATETIME_FORMAT)
        except Exception:
            # NOTE: this also gives a pass to any malformed datetime strings!
            # It's necessary, because otherwise it throws an error for any
            # non-datetime string (e.g. "test").
            # MAKE SURE that all datetime strings you want to deserialize
            # have been serialized with json_serial function. For example,
            # they always have to have timezone set!
            pass
    return json_dict


def grpcdataToDict(grpcdata):
    """
    Converts bytes data received over gRPC to a Dictionary object.
    """
    datastr = grpcdata.decode('utf-8')
    dic = json.loads(datastr, object_hook=json_deserial)
    return dic
