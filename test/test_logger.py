import logging
from src.utils.Logger import setup_logger


def log_all_levels(log):
    log.debug('debug message')
    log.info('info message')
    log.warning('warning message')
    log.error('error message')
    log.critical('critical message')


def capture_logged_messages(capfd):
    _, err = capfd.readouterr()
    captured_logs_list = err.split('\n')
    return captured_logs_list


def test_log_all_levels(capfd,
                        log=setup_logger(
                            logging.getLogger('test_DEBUG'),
                            handler=logging.StreamHandler(),
                            level='DEBUG')):
    log_all_levels(log)
    captured_logs_list = capture_logged_messages(capfd)
    assert any('DEBUG' in message for message in captured_logs_list)


def test_log_INFO_up(capfd,
                     log=setup_logger(
                         logging.getLogger('test_INFO'),
                         handler=logging.StreamHandler(),
                         level='INFO')):
    log_all_levels(log)
    captured_logs_list = capture_logged_messages(capfd)
    assert not any('DEBUG' in message for message in captured_logs_list)


def test_log_WARNING_up(capfd,
                        log=setup_logger(
                            logging.getLogger('test_WARNING'),
                            handler=logging.StreamHandler(),
                            level='WARNING')):
    log_all_levels(log)
    captured_logs_list = capture_logged_messages(capfd)
    assert not any('INFO' in message for message in captured_logs_list)


def test_log_ERROR_up(capfd,
                      log=setup_logger(
                          logging.getLogger('test_ERROR'),
                          handler=logging.StreamHandler(),
                          level='ERROR')):
    log_all_levels(log)
    captured_logs_list = capture_logged_messages(capfd)
    assert not any('WARNING' in message for message in captured_logs_list)


def test_log_CRITICAL_up(capfd,
                         log=setup_logger(
                             logging.getLogger('test_CRITICAL'),
                             handler=logging.StreamHandler(),
                             level='CRITICAL')):
    log_all_levels(log)
    captured_logs_list = capture_logged_messages(capfd)
    assert not any('ERROR' in message for message in captured_logs_list)
