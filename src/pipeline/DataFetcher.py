from src.components import TermSplitter
from src.utils.measures import measure_time
import pandas as pd
import logging
from prometheus_client import Summary, Gauge
import os
from datetime import timedelta, datetime, timezone
import pickle

logger = logging.getLogger('root')
current_env_file = os.path.splitext(os.path.basename(__file__))[0]
prometheus_summary_conference_counter = Summary(
    current_env_file + '_conference_counter_time', 'Summary time and count')

fetched_data_prometheus = Gauge('Fetched_data_prometheus', 'appIDs',
                                ['app_id', 'metric_name'])

# NOTE: TEMPORARY
# TODO: Remove once CRS has enough effects data
BACKFILL_START = datetime(2018, 10, 26, tzinfo=timezone.utc)
BACKFILL_END = datetime(2019, 2, 20, tzinfo=timezone.utc)
with open('backfill.pickle', 'rb') as handle:
    BACKFILL_DATA = pickle.load(handle)
DATE_FORMAT = '%Y%m%d'


class DataFetcher(TermSplitter):
    """
    A DataFetcher object fetches required data for the models for all appIDs.
    """

    def __init__(self, metrics, appIDs, crsClient):
        """
        parameters:
            appIDs: list, appIDs for which to fetch metrics
            crsClient: CrsClient, connection to CRS
        """
        logger.debug('Creating DataFetcher object')
        self._app_ids = appIDs
        self._metrics = metrics
        self._crs_client = crsClient
        TermSplitter.__init__(self, process=self._receive)
        logger.debug('DataFetcher created')

    def _is_within_backfill_timerange(self, date_package):
        return ((date_package.start >= BACKFILL_START)
                and (date_package.start <= BACKFILL_END))

    def _handle_prometheus(self, data):
        for metric_name, metric_data in data.items():
            [
                self._send_prometheus_point(app_id, app_data, metric_name)
                for metric_name, metric_data in data.items()
                for app_id, app_data in metric_data.items()
            ]

    def _send_prometheus_point(self, app_id, app_data, metric_name):
        fetched_data_prometheus.labels(app_id, metric_name).set(
            app_data.reset_index(drop=True).value[0])

    @prometheus_summary_conference_counter.time()
    @measure_time
    def _receive(self, date_package):
        """
        Receives a DatePackage and generates metrics
        input:
            date_package: DatePackage with set start and end
        output:
            dict, entry for each appID:
                time: Pandas datetime, end of time frame
                value: int, number of conferences
        """
        aggregated = {}
        for appID in self._app_ids:
            metrics = self._crs_client.Aggregate(
                appID,
                date_package.start,
                # NOTE CRS includes all full DAYS into the query for which
                # there is a datetime including (not excluding) end.
                # So we have to substract a bit to get desired results.
                date_package.end - timedelta(milliseconds=1))
            logger.debug('Data fetched: {} - {}'.format(appID, metrics))
            aggregated[appID] = metrics

        # Filter only metrics we support for AID
        result = {
            m: {
                appID: pd.DataFrame({
                    'time': [date_package.end],
                    'value': aggregated[appID][m]
                })
                for appID in self.app_ids
            }
            for m in self.metrics
        }
        logger.debug('Values got from CRS {}'.format(result))
        if self._is_within_backfill_timerange(date_package):
            update_val = BACKFILL_DATA[date_package.start.strftime(
                DATE_FORMAT)]
            result.update(update_val)
            logger.debug('Used backfill values {}'.format(update_val))
        logger.debug('Final metric values {}'.format(result))
        self._handle_prometheus(result)
        return result

    @property
    def app_ids(self):
        return self._app_ids

    @property
    def metrics(self):
        return self._metrics
