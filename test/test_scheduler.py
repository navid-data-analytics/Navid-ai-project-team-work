"""Tests scheduler class."""
from src.components import Collector
from src.pipeline.Scheduler import Scheduler
import src.utils.constants as constants
import datetime
import pytest


def test_trigger_interval_under_threshold_raises_value_error(
        time_interval=0.5 * constants.SECOND):
    with pytest.raises(ValueError):
        Scheduler(trigger_interval=time_interval)


def send_initial_load(scheduler, output_collector):
    scheduler.output.connect(output_collector.input)
    scheduler.start()
    scheduler.shutdown()
    output_collector.input(None)


def prep_scheduler(scheduler):
    scheduler._scheduler.add_job(
        func=lambda: scheduler._notify(),
        trigger='interval',
        seconds=scheduler._trigger_interval,
        next_run_time=datetime.datetime.utcnow(),
        misfire_grace_time=constants.MINUTE)


def test_init_buffer_amount(trigger_interval=2 * constants.YEAR,
                            init_buffer_size=5):
    output_collector = Collector()
    scheduler = Scheduler(trigger_interval, init_buffer_size)
    prep_scheduler(scheduler)
    send_initial_load(scheduler, output_collector)
    assert len(output_collector._items) == init_buffer_size


def test_initial_buffer_start_time(trigger_interval=2 * constants.YEAR,
                                   init_buffer_size=5):
    scheduler = Scheduler(trigger_interval, init_buffer_size)
    scheduler._set_initial_start_time()
    catch_up_time = scheduler._get_current_time()
    start_time = scheduler._start_time
    assert start_time != catch_up_time, 'Start should be 5 intervals before!'


def test_initial_buffer_end_time(trigger_interval=2 * constants.YEAR,
                                 init_buffer_size=5):
    scheduler = Scheduler(trigger_interval, init_buffer_size)
    scheduler._set_initial_start_time()
    prep_scheduler(scheduler)
    catch_up_time = scheduler._get_current_time()
    scheduler.start()
    scheduler.shutdown()
    last_sent_end_time = scheduler._end_time
    assert last_sent_end_time == catch_up_time, '''Last end time should be
 equal to catch_up_time!'''


def check_truncation(scheduler, trigger_interval):
    time_unit = get_appropriate_unit(trigger_interval)
    start_datetime = datetime.datetime.fromtimestamp(
                        scheduler._start_time / constants.CONVERT_KILO)
    assert getattr(start_datetime, time_unit) == 0


def get_appropriate_unit(trigger_interval):
    if trigger_interval < constants.MINUTE:
        time_unit = 'microsecond'
    elif trigger_interval < constants.HOUR:
        time_unit = 'second'
    elif trigger_interval < constants.DAY:
        time_unit = 'minute'
    else:
        time_unit = 'hour'
    return time_unit


def test_time_truncation_seconds(trigger_interval=2 * constants.SECOND,
                                 init_buffer_size=5):
    output_collector = Collector()
    scheduler = Scheduler(trigger_interval, init_buffer_size)
    send_initial_load(scheduler, output_collector)
    check_truncation(scheduler, trigger_interval)


def test_time_truncation_minutes(trigger_interval=2 * constants.MINUTE,
                                 init_buffer_size=5):
    output_collector = Collector()
    scheduler = Scheduler(trigger_interval, init_buffer_size)
    send_initial_load(scheduler, output_collector)
    check_truncation(scheduler, trigger_interval)


def test_time_truncation_hours(trigger_interval=2 * constants.HOUR,
                               init_buffer_size=5):
    output_collector = Collector()
    scheduler = Scheduler(trigger_interval, init_buffer_size)
    send_initial_load(scheduler, output_collector)
    check_truncation(scheduler, trigger_interval)


def test_time_truncation_days(trigger_interval=2 * constants.DAY,
                              init_buffer_size=5):
    output_collector = Collector()
    scheduler = Scheduler(trigger_interval, init_buffer_size)
    send_initial_load(scheduler, output_collector)
    check_truncation(scheduler, trigger_interval)
