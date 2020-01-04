import time
import datetime as dt


def to_dtutc(current_time):
    """
    Convert the timeStamp type to datetime with local UTC.

    Arguments:
    - current_time: Pandas timeStamp, the time to be converted

    Returns:
    - datetime.datetime, converted time format
    """
    return current_time.to_pydatetime()


def to_unix_timestamp(current_time):
    """
    Convert the timeStamp type to datetime with local UTC.

    Arguments:
    - current_time: datetime object to be converted

    Returns:
    - float, unix timestamp
    """
    return time.mktime(current_time.timetuple())


def get_time_dict(base_time,
                  previous=(None, None),
                  current=(None, None),
                  future=(None, None)):
    """
    Get a period based time dict.

    Arguments:
    - base_time: datetime object, base time to work off
    - remaining: tuple of two ints, date relative to the base_time in days, in
                 ('start', 'end') format. None on either field will result in
                 the field not being present in result

    Returns:
    - dict, {'period1': {'start': base_time + start_delta,
                        'end': base_time + end_delta}
            (...)
            }
    """
    key_names = ['previous_period', 'current_period', 'future_period']
    delta_tuples = [previous, current, future]
    return {key: get_period_dict(base_time, delta_tuple) for key, delta_tuple
            in zip(key_names, delta_tuples) if delta_tuple != (None, None)}


def get_period_dict(base_time, delta_tuple):
    return {key: base_time + dt.timedelta(days=value) for key, value in
            zip(['start', 'end'], delta_tuple) if value is not None}
