"""This file contains Scheduler class and every dependencies for it."""
from apscheduler.schedulers.background import BackgroundScheduler
from src.components import Transmitter, DatePackage
from src.utils.measures import measure_time
import src.utils.constants as constants
import datetime
import logging


logger = logging.getLogger('root')


class Scheduler(Transmitter):
    """
    Scheduler is responsible for initiating the Anomaly detection pipeline.

    A Scheduler object encapsulates all the logic necessary for managing jobs
    between Scheduler and DataFetcher object. The object contains clock,
    current trigger info, initial_load and minimal trigger threshold. The clock
    is a wrapped APScheduler library BackgroundScheduler object, with run,
    shutdown methods implemented. The trigger interval value has to be BIGGER
    than MIN_INTERVAL_THRESHOLD value which is set to 1 second.

    Output:
    - DatePackage with time range filled (start, end), no data beyond time.

    The workflow should look like this:
    - Scheduler is created with trigger interval and initial load specified.
    - Scheduler sends the initial buffer of size equal to trigger_interval,
      initial_load steps back.
    - After trigger interval the Scheduler sends the DatePackage with correctly
      time range.
    """

    MIN_INTERVAL_THRESHOLD = constants.SECOND
    LAUNCH_HOUR = 12

    def __init__(self,
                 trigger_interval,
                 initial_load=0,
                 *args,
                 **kwargs):
        """
        Construct a new Scheduler instance.

        Arguments:
        - trigger_interval: integer, time after which scheduler sends
          DatePackage, also governs the size of time_range of the DatePackage.
        - initial_load: integer, amount of DatePackage objects Scheduler should
          send on start, by default set to 0.

        raises: ValueError
        - in case trigger_interval is below the MIN_INTERVAL_THRESHOLD.
        """
        if trigger_interval < self.MIN_INTERVAL_THRESHOLD:
            raise ValueError("""The trigger interval cannot be
                              lower than the threshold value""")

        logger.debug('Creating Scheduler object')
        Transmitter.__init__(self, process=lambda: None)
        logger.debug('Storing trigger interval and initial load info')
        self._trigger_interval = trigger_interval
        self._initial_load = initial_load
        self._set_up_scheduler(*args, **kwargs)
        logger.debug('Scheduler created')

    @property
    def trigger_interval(self):
        return self._trigger_interval

    def start(self, *args, **kwargs):
        """Start the timer in the scheduler object."""
        self._set_initial_start_time()
        self._trigger_initial_run()
        logger.debug('Starting clock')
        self._scheduler.start(*args, **kwargs)

    def shutdown(self):
        """Shutdown the scheduler object."""
        logger.debug('Shutting down clock')
        self._scheduler.shutdown()

    def _set_up_scheduler(self, *args, **kwargs):
        logger.debug('Setting up Scheduler backend')
        self._scheduler = BackgroundScheduler(*args, **kwargs)
        first_run_time = self._set_launch_time()
        self._scheduler.add_job(
            func=lambda: self._notify(),
            trigger='interval',
            seconds=self._trigger_interval,
            next_run_time=first_run_time,
            misfire_grace_time=constants.MINUTE,
            *args,
            **kwargs)
        logger.debug('Done setting up Scheduler backend')

    def _set_launch_time(self):
        """Set launch time."""
        logger.debug('Setting up launch time')
        utc_now = datetime.datetime.now(datetime.timezone.utc)
        logger.debug('Checking if launch should happen today')
        if utc_now.hour >= self.LAUNCH_HOUR:
            logger.info(
                'It is past {}, next run time set to tomorrow.'
                .format(self.LAUNCH_HOUR))
            utc_now += datetime.timedelta(days=1)
        run_time = utc_now.replace(hour=self.LAUNCH_HOUR,
                                   minute=0,
                                   second=0,
                                   microsecond=0)
        logger.debug(
            'Launch time set to {}'.format(
                run_time.strftime('%H:%M:%S %d-%m-%Y')))
        return run_time

    @measure_time
    def _notify(self):
        current_time = self._get_current_time()
        while self._start_time < current_time:
            self._end_time = self._start_time + self._trigger_interval \
                * constants.CONVERT_KILO
            self._send_date_package(self._start_time, self._end_time)
            self._start_time = self._end_time

    def _get_current_time(self):
        current_time = datetime.datetime.now(datetime.timezone.utc)
        current_time = self._truncate_interval_residue(current_time)
        return current_time

    def _truncate_interval_residue(self, time):
        logger.debug("Truncating current time")
        if self._trigger_interval < constants.MINUTE:
            truncation_dict = {'microsecond': 0}
        elif self._trigger_interval < constants.HOUR:
            truncation_dict = {'second': 0, 'microsecond': 0}
        elif self._trigger_interval < constants.DAY:
            truncation_dict = {'minute': 0, 'second': 0, 'microsecond': 0}
        else:
            truncation_dict = {
                'hour': 0,
                'minute': 0,
                'second': 0,
                'microsecond': 0
            }
        time = time.replace(**truncation_dict).timestamp() \
            * constants.CONVERT_KILO
        return time

    def _set_initial_start_time(self):
        logger.debug("Setting initial start time")
        current_time = self._get_current_time()
        self._start_time = current_time - (self._initial_load *
                                           self._trigger_interval *
                                           constants.CONVERT_KILO)
        logger.debug("Initial start time set to: {}".format(self._start_time))

    def _send_date_package(self, start, end):
        logger.debug(
            "Creating date package object with start end: {} {}".format(
                start, end))
        date_package = DatePackage(start, end)
        logger.debug("Send DatePackage: {}".format(date_package))
        self.output.transmit(date_package)

    def _trigger_initial_run(self):
        # NOTE: This function allows pipeline to run regardless
        #       of the launch hour - which tells at what time
        #       the pipeline is triggered daily. In other cases
        #       than initial the pipeline will be launched always
        #       at the same time, specified by LAUNCH_HOUR.
        logger.debug('Run the Scheduler Once to process data up to today.')
        self._notify()
        logger.debug('Initial trigger finished.')
