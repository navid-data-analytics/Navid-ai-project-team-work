import grpc
import time
import logging

logger = logging.getLogger('root')


class DataServiceError(BaseException):
    def __init__(self, endpoint_name, grpc_error, additional_info=""):
        super(DataServiceError,
              self).__init__('data_service.{}: {} ({}) {}'.format(
                  endpoint_name, grpc_error.details(),
                  grpc_error.code().name, additional_info))


class ConnectionClient(object):
    """
    A gRPC client. It handles the connection level, i.e.
    establishment and reliability.
    """
    def __init__(self, address, connection_timeout, max_retries):
        """
        input:
            address: String "<ip addr or name>:<port>" of the server
            connection_timeout: int, seconds, timeout of connection to server
            max_retries: int, how often to retry if fails,
                0 means retrying indefinitely
        """
        self._address = address
        self._connection_timeout = connection_timeout
        self._connection = None
        self._max_retries = max_retries

    @property
    def connection(self):
        """
        returns: a gRPC connection object
        """
        if self._connection is None:
            self._connection = grpc.insecure_channel(self._address)
            # This tests we can connect to the service within given timeout,
            # essentially making sure the service is running, before we start
            # sending requests to it.
            # If we didn't check this, the requests would take forever.
            # If the timeout occurs, `grpc.FutureTimeoutError` will be thrown.
            grpc.channel_ready_future(
                self._connection).result(timeout=self._connection_timeout)
        return self._connection

    def getService(self, stub):
        """
        creates the gRPC service.
        input:
            stub: a gRPC stub method
        returns:
            gRPC service object
        """
        service = None
        while service is None:
            try:
                service = stub(self.connection)
            except (grpc.FutureTimeoutError):
                # FutureTimeoutError stops the waiting for channel ready.
                # But we can ignore it, the next error will be UNAVAILABLE
                continue
        return service

    def send(self, method, request, name, reliable=False):
        """
        sends a gRPC request either reliably or unreliably.
            reliable: retry if service is unavailable
            unreliable: try only once and fail if unavailable
        input:
            method: gRPC service method
            request: gRPC request method corresponding to the method
            name: String, for debugging
            reliable: Boolean (see above)
        returns:
            a tuple of (value, Exception), value is method specific
            value can be None, if no response generated or error occured
            error is either an Exception, or None if no error occured
        """
        delay = 1
        retries = 0
        while True:
            retries += 1
            try:
                return method(request), None
            except (grpc.RpcError) as e:
                if reliable and hasattr(e, 'code'):
                    code = e.code()
                    if code == grpc.StatusCode.UNAVAILABLE:
                        if self._max_retries > 0 and \
                           retries > self._max_retries:
                            info = "too many retries (max={})"\
                                     .format(self._max_retries)
                            err = DataServiceError(name, e, info)
                            logger.error(err)
                            return None, err
                        logger.warn(
                            DataServiceError(name, e,
                                             "retry in %ds" % (delay)))
                        time.sleep(delay)
                        if delay < 30:
                            delay += delay
                        continue

                err = DataServiceError(name, e)
                logger.error(err)
                return None, err
