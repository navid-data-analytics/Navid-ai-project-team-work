import datetime
import src.utils.constants as constants


class DatePackage:
    def __init__(self, start, end, data=None):
        """
        inputs:
            start: float, timestamp in ms
                OR datetime
            end: float, timestamp in ms
                OR datetime
        """
        if isinstance(start, datetime.datetime):
            self._start = start
        else:
            self._start = self._timestamp_ms_to_datetime(start)

        if isinstance(end, datetime.datetime):
            self._end = end
        else:
            self._end = self._timestamp_ms_to_datetime(end)

        self.set_data(data)

    @property
    def start(self):
        return self._start

    @property
    def end(self):
        return self._end

    @property
    def data(self):
        return self._data

    def set_data(self, data):
        self._data = data

    def __str__(self):
        return "DatePackage(%s, %s, %s)" % (
            self._start, self._end, type(self._data))

    def _timestamp_ms_to_datetime(self, timestamp):
        return datetime.datetime.fromtimestamp(
            timestamp / constants.CONVERT_KILO, datetime.timezone.utc)
