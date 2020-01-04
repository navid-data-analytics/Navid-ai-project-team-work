"""Tests DateQueue class."""
from src.components import DatePackage
from src.pipeline import DateQueue
import src.utils.constants as constants
import time


def test_datequeue_queues_job(expected_outcome=False):
    datequeue = DateQueue()
    datequeue.input(None)  # stop the listening thread

    datequeue.input(DatePackage(start=0,
                                end=time.time() * constants.CONVERT_KILO))
    result = datequeue.requests.empty()
    assert result == expected_outcome
    datequeue.shutdown()


def test_datequeue_processes_jobs(expected_outcome=True):
    datequeue = DateQueue()

    datequeue.input(DatePackage(start=0,
                                end=time.time() * constants.CONVERT_KILO))
    datequeue.input(DatePackage(start=0,
                                end=time.time() * constants.CONVERT_KILO))
    time.sleep(0.5)  # wait for other thread to work
    result = datequeue.requests.empty()
    assert result == expected_outcome
    datequeue.shutdown()
