"""File contains all logic necessary for capturing logs."""
from raven.handlers.logging import SentryHandler
from raven.conf import setup_logging
from raven import Client
import logging

logger = logging.getLogger('root')


def setup_logger(logger, handler, level):
    """
    Set up the logger.

    Arguments:
    - logger: The logger used for the whole execution. Usually it should
      be logging.getLogger('root').
    - handler: Handler determines type of output logging and style.
    - level: Specifies the minimal level of logs that should be captured:
      DEBUG, INFO, WARNING, ERROR, CRITICAL. Picking one captures all levels
      above it.
    """
    fmt = '[%(asctime)s.%(msecs)03d] (%(threadName)s) \
%(filename)s %(levelname)s: %(message)s'
    fmt_date = '%Y-%m-%d %T'
    handler.setFormatter(logging.Formatter(fmt, fmt_date))
    logger.setLevel(level)
    logger.addHandler(handler)
    return logger


def setup_sentry(sentry_credentials, release):
    """
    Set up handler for sentry.

    With this handler working all logs with level ERROR or higher
    (ERROR/CRITICAL) are sent directly to sentry.io, all others all left
    as they are. Steps to set up handler are according to the official docs:
    https://github.com/getsentry/raven-python.

    returns:
    - sentry_handler: object of SentryHandler class.
    """
    logger.info('Setting up sentry handler and client.')
    if sentry_credentials:
        client = Client(sentry_credentials)
        sentry_handler = SentryHandler(sentry_credentials,
                                       release=release)
        sentry_handler.setLevel(logging.ERROR)
        setup_logging(sentry_handler)
        logger.info('Sentry handler set up correctly, returning client')
        return client
    else:
        logger.error("Sentry credentials are incorrect.")
