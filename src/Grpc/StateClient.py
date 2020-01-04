from src.Grpc.AidServiceClient import AidServiceClient
from collections import defaultdict
import datetime
import logging
import itertools

logger = logging.getLogger('root')

DEFAULT_DT = datetime.datetime(1971, 1, 1, tzinfo=datetime.timezone.utc)
PREVIOUS_LATEST_DATE_TAG = 'latest_date'
STATE_LATEST_DATE_TAG = 'latest_dates'
STATE_LATEST_DATE_DATA_ENTRY = 'state'


class StateClient(AidServiceClient):
    """
    Intermediate class between AidServiceClient and MessageClient to handle
    pipeline state loading and usage.
    """

    @staticmethod
    def _get_greater(a, b):
        return a if a > b else b

    @staticmethod
    def _parse_dt(date):
        return datetime.datetime.strptime(
            date, '%d-%m-%Y').replace(tzinfo=datetime.timezone.utc)

    def __init__(self, address, flags, load_state):
        super(StateClient, self).__init__(address)
        self._default_date = DEFAULT_DT
        self._set_flags(flags)
        self._set_latest_dates(load_state)
        self._handle_manual_update()

    def _set_latest_dates(self, load_state):
        """Set latest date to be used for message suppression.

        inputs:
            - load_state: bool, whether state should be loaded
        """
        if not load_state:
            self._latest_dates = defaultdict(dict)
            return

        self._latest_dates = defaultdict(dict, self._get_dates_from_state())
        logger.debug('Latest state: {}'.format(self._latest_dates))

    def _handle_if_malformed(self, state):
        if STATE_LATEST_DATE_DATA_ENTRY in state:
            return state
        logger.error('State is malformed')
        return {}

    def _get_dates_from_state(self):
        """Get state dict from state grpc message.

        The expected structure of received message is a dict, where state is
        under 'state' key. The state itself is a two level dictionary, with
        first level keys being appIDs and second level having
        'message_type: datetime' pairs. If instead the state is a single
        datetime object, this date will be used for all appIDs and
        message_types

        returns:
            - datetime, latest date
        """
        state = self.GetState(STATE_LATEST_DATE_TAG)
        logger.info('State received: {}'.format(state))

        if state is None:
            logger.info('No state received, attempting to use previous ver.')
            state = self.GetState(PREVIOUS_LATEST_DATE_TAG)
            if state is None:
                logger.info('No usable state found')
                return {}
            else:
                logger.info('Using previous version state')
                state = self._handle_if_malformed(state)
                result = state[STATE_LATEST_DATE_DATA_ENTRY]
                logger.info(
                    'Setting universal suppression date to {}'.format(result))
                self._default_date = result
                return {}

        state = self._handle_if_malformed(state)
        result = self._del_faulty_state_entries(
            state[STATE_LATEST_DATE_DATA_ENTRY])
        logger.info('Loaded state: {}'.format(result))
        return result

    def _del_faulty_state_entries(self, state):
        for app_id in state.keys():
            for message_type in list(state[app_id].keys()):
                date = state[app_id][message_type]
                if not isinstance(date, datetime.datetime):
                    logger.error('State for {} {} is malformed'.format(
                        app_id, message_type))
                    del state[app_id][message_type]
        return state

    def _handle_suppression(self, dt, appID, type):
        appID = str(appID)
        appid_state = self._latest_dates[appID]
        type_date = appid_state.get(type, self._default_date)
        if type_date == self._default_date:
            self._latest_dates[appID][type] = dt
        if dt <= type_date:
            if self._check_unsuppress(appID, dt, type):
                logger.debug(
                    'Message {} app_id {} not suppressed, {} == {}'.format(
                        type, appID, dt, type_date))
                return True
            logger.info('Message {} app_id {} suppressed, {} < {}'.format(
                type, appID, dt, type_date))
            return False
        logger.debug('Message {} app_id {} not suppressed, {} > {}'.format(
            type, appID, dt, type_date))
        return True

    def _CreateMessageConsiderState(self, dt, appID, type, version, data):
        """
        Wraps ProtocolBuffer's _CreateMessage, considering the date to decide
        whether the message should be sent or suppressed.

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
        if self._handle_suppression(dt, appID, type):
            return self._CreateMessage(dt, appID, type, version, data)
        return None

    def save_dates(self, dt):
        """Save date as current state.

        input:
            dt: datetime, date to be saved
        returns:
            Exception, None if no error
        """
        for appid in self._latest_dates.keys():
            for type in self._latest_dates[appid].keys():
                if dt > self._latest_dates[appid][type]:
                    self._latest_dates[appid][type] = dt
        logger.info('Save state: {}'.format(self._latest_dates))
        return self.SaveState(
            STATE_LATEST_DATE_TAG,
            {STATE_LATEST_DATE_DATA_ENTRY: self._latest_dates})

    def _check_unsuppress(self, appID, dt, msg_type):
        flags_are_empty = self.flags.get('unsuppress', None) in ({}, None)
        logger.debug("dt: {} appid: {} flags are empty?: {}".format(
            dt, appID, flags_are_empty))
        if flags_are_empty:
            return False
        try:
            unsuppress_datetimes = list(
                map(self._parse_dt, self.flags['unsuppress'][appID][msg_type]))
            if dt in unsuppress_datetimes:
                return True
            return False
        except KeyError:
            return False

    def _handle_manual_update(self):
        flags_are_empty = self.flags.get('date', None) in ({}, None)
        if flags_are_empty:
            return
        self._update_dates()

    def _update_dates(self):
        appIDs, message_types, dates = self._parse_manual_date()
        logger.info(
            'Manually set state for appid {} message type {} date {}'.format(
                appIDs, message_types, dates))
        self._update_latest_dates(appIDs, message_types, dates)

    def _parse_manual_date(self):
        appIDs = list(self.flags['date'].keys())
        message_types = list(
            itertools.chain.from_iterable(
                [i.keys() for i in self.flags['date'].values()]))
        datetime_dates = self._parse_flag_dates(key='date')
        return appIDs, message_types, datetime_dates

    def _update_latest_dates(self, appIDs, message_types, dates):
        for appID in appIDs:
            for idx, message_type in enumerate(message_types):
                self._latest_dates[appID][message_type] = dates[idx]

    def _parse_flag_dates(self, key):
        dates = []
        for i in self.flags[key].values():
            dates.extend(list(itertools.chain.from_iterable(i.values())))
        datetime_dates = [self._parse_dt(date) for date in dates]
        return datetime_dates

    def _set_flags(self, flags):
        keys = flags.keys()
        self._flags = flags
        logger.debug('Flags set to {}'.format(self.flags))
        if 'date' not in keys:
            self._flags['date'] = None
        if 'unsuppress' not in keys:
            self._flags['unsuppress'] = None

    @property
    def flags(self):
        return self._flags
