"""Fetches the data from the cluster every time the scheduler tells it to."""
from src.components.transceiver import Transceiver
from threading import Thread
from queue import Queue
import logging

logger = logging.getLogger('root')


class DateQueue(Transceiver):
    """
    DateQueue buffers DatePackage objects.

    Input:
    - DatePackage object with time range already specified.

    Output:
    - DatePackage with the same time range.
    """

    def __init__(self):
        """
        Construct a new DateQueue instance.

        Fetcher uses DatePackage object's start and end parameters
        (given on input) to extract the data from MongoDB.

        params:
        - path: Path to the csv_file data should be load from
        """
        Transceiver.__init__(self, process=lambda: None)
        logger.debug('Creating blocking queue')
        self._requests = Queue()
        logger.debug('Creating fetching thread')
        self._fetcher_thread = Thread(target=self._listen_queue)
        logger.debug('Starting fetching thread')
        self._fetcher_thread.start()

    def show_queue(self, q):
        with q.mutex:
            return list(q.queue)

    @property
    def thread(self):
        return self._fetcher_thread

    @property
    def requests(self):
        return self._requests

    def input(self, date_package):
        """
        Push received date_package object into request queue on input.

        params:
        - date_package: DatePackage object containing start and end times.
        """
        logger.debug('Received date package on input, inserting to requests')
        self._push_request(date_package)
        logger.debug('Done inserting date package to requests queue')

    def shutdown(self):
        """Stop the Queue thread from execution."""
        logger.debug('Shutting down DateQueue')
        self.input(None)

    def _push_request(self, date_package):
        logger.debug('Contents of queue before adding job: {}'.format(
            self.show_queue(self._requests)))
        self._requests.put(date_package, block=True)
        logger.debug(
            'Contents of queue after adding job: {}'.format(
             self.show_queue(self._requests)))

    def _listen_queue(self):
        while True:
            logger.debug('Starting _listen_queue() iteration.')
            logger.debug(
                'Contents of queue at the beginning of listen_queue: {}'
                .format(self.show_queue(self._requests)))
            logger.debug('Get current job from queue {}'.format(
                self.show_queue(self._requests)))
            current_job = self._requests.get(block=True)
            logger.debug('Fetching started')
            if current_job is None:
                logger.debug('Got None, exiting queue.')
                break
            logger.info('Current time range: {}'.format(current_job))
            self._output.transmit(current_job)

            logger.debug('Calling task_done()')
            self._requests.task_done()
            logger.debug('Set current_job to None')
            current_job = None
            logger.debug('Set current_job to None successfully.')
            logger.debug(
                'Contents of queue at the end of listen_queue: {}'
                .format(self.show_queue(self._requests)))
            logger.debug('Finishing _listen_queue() iteration')
        logger.debug('Terminating fetching queue')
