"""The AID main program."""

import ast
import signal
import logging
import traceback
import threading
import argparse
from prometheus_client import generate_latest, start_http_server
from src.utils import setup_sentry
from src.pipeline import Pipeline, EnvConfig
from src.Grpc.MessageClient import MessageClient
from src.Grpc.CrsClient import CrsClient

env = EnvConfig()
logger = logging.getLogger('root')

parser = argparse.ArgumentParser()
parser.add_argument(
    '--unsuppress',
    type=str,
    help='''
Unsuppress Grpc State saving for a passed date.

The format of unsuppression should be as follows:

--unsuppress {appID: {message_type: date} }

raw example:

--unsuppress {7233: {'MidtermRttFluctuationStabilized': "20-01-2018"}}

For example above, appid 7233 for message type
MidtermRttFluctuationStabilized
will be unsuppressed for 20th of January, 2018.
Make sure the message conflicts do not occur by deleting
the messages from database beforehand.
''')
parser.add_argument(
    '--manual_date',
    type=str,
    default='{}',
    help='''
Update manually the date for suppression mechanism in gRPC client.

The format of manual update should be as follows:

    --manual_date {appID: {message_type: date} }

raw example:

    --manual_date {7233: {'MidtermTrend15daysUp': "20-01-2018"}}

For example above, appid 7233 for message type MidtermTrend15daysUp
last suppressed date will be updated to 20th of January, 2018.
''')

args = parser.parse_args()

logger.info("Flags passed:\nmanual_date: {}\nunsuppress: {}".format(
    args.manual_date, args.unsuppress))
flags = {
    'unsuppress': ast.literal_eval(args.unsuppress),
    'date': ast.literal_eval(args.manual_date)
}


class Thread_Prometheus(threading.Thread):
    def run(self):
        start_http_server(env.port_prometheus)
        return generate_latest()


if __name__ == '__main__':

    thread_prometheus = Thread_Prometheus()
    thread_prometheus.start()

    sentry_client = setup_sentry(env.sentry_credentials, env.version)

    AidServiceConnection = MessageClient(
        env.aid_service_grpc_address, flags, load_state=True)
    CrsClient = CrsClient(env.crs_grpc_address)

    pipeline = Pipeline(env, AidServiceConnection, CrsClient)
    try:
        pipeline.start()
        signal.pause()
    except Exception as e:
        logger.critical('Internal crash of the pipeline happened!\n'
                        'Exception: {}'.format(e) + traceback.format_exc())
        # locally we do not connect to sentry
        if sentry_client:
            sentry_client.captureException()
    finally:
        pipeline.stop()
