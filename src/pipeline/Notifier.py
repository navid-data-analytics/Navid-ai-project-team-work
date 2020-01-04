"""Notify sentry.io in case score was anomalous."""
from raven import Client
from src.components import Receiver
from src.utils import measure_time
import logging
import os

logger = logging.getLogger('root')

current_env_file = os.path.splitext(os.path.basename(__file__))[0]


class Notifier(Receiver):
    """
    A Notifier class is responsible connecting to sentry.io.

    Notifier object should be constructed with sentry credentials as input
    to the constructor.
    """

    def __init__(self,
                 client_params=None,
                 release=None,
                 message_client_address=None,
                 aid_service_connection=None):
        self._client = None
        Receiver.__init__(self, process=self._run)
        """Construct a new Notifier instance."""
        logger.debug('Creating Notifier object')
        if client_params:
            logger.info('Connecting to sentry: client: {}, release: {}'.format(
                client_params, release))
            self._client = Client(client_params, release=release)
        self._aid_service_connection = aid_service_connection
        self._message = None
        logger.debug('Notifier created')

    @measure_time
    def _run(self, input):

        logger.debug('Input of Notifier: {}'.format(input))
        date = self._get_date(input)
        # self._handle_oq_complex_message(input)
        self._save_state(date)

    def _get_date(self, input):
        metrics = list(input.keys())
        appids = list(input[metrics[0]].keys())
        detections = list(input[metrics[0]][appids[0]].keys())
        return input[metrics[0]][appids[0]][detections[0]]['time_model']

    def _handle_oq_complex_message(self, input):
        for term in ('shortterm', 'midterm'):
            oq_data = input[term]['complex_dm_output'].get(
                'objective_quality_v35_average', None)
            if oq_data is not None:
                for appid in oq_data.keys():
                    self._send_message_per_app(oq_data[appid])

    def _send_complex_message(self, input):
        if input['send'] is False:
            return
        direction = input['direction']
        trend_type = input['type']
        reasons = self._determine_reasons(input)
        if len(reasons) < 2:
            return
        self._notify(trend_type, direction, reasons, input['message_details'])

    def _determine_reasons(self, input):
        reasons = [
            ''.join(reason.split("_")[:-1]) for reason in input['reasons']
            if 'midterm' in reason
        ]

        result = ['objective_quality_v35_average'] + reasons
        return tuple(sorted(result))

    def _send_message_per_app(self, oq_data):
        for appID in oq_data.keys():
            input_for_appid = oq_data[appID]
            self._send_complex_message(input_for_appid)

    def _notify(self, trend_type, direction, reasons, message_details):
        logger.debug('Sending complex message!')
        logger.info('REASONS: {}, deets {}'.format(reasons, message_details))
        self._message = self._aid_service_connection.messages[
            'objective_quality_v35_average']['trend'][direction][trend_type][
                reasons](*message_details)

    def _save_state(self, date):
        if self._aid_service_connection:
            logger.debug('Processed date {}, saving state'.format(date))
            self._aid_service_connection.save_dates(date)
