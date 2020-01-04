# pylama:ignore=E731
import numpy as np

is_odd = lambda x: x % 2 == 1


def running_mean_fast(signal, window_size):
    """
        Rapidly computes running mean of a signal via convolution.
        Args:
            - signal (1d array): original signal
            - window_size (int): window size for running mean calculation
        Returns:
            - 1d array: running mean of the signal (array shorter by
                        window_size)
    """
    return np.convolve(signal, np.ones((window_size,))/window_size,
                       mode='valid')


def running_fn(signal, fn, window_size, subtract_mean=True):
    """
    Compute a running function of a signal.

    Args:
        - signal (1d array): original signal
        - fn (function): the function that will be ran over windowed signal
        - window_size (int): window size, used for both running mean and
                             running function
        - subtract_mean (bool): whether the mean of the signal should be
                                subtracted before computing the running
                                function
    Returns:
        - 1d array: result of running function over the original signal
    """
    half_window = int(window_size/2)
    endpoint = signal.shape[0] - half_window
    if subtract_mean:
        mean = running_mean_fast(signal, window_size)
        if is_odd(window_size):
            signal[half_window:endpoint] = (
                signal[half_window:endpoint] - mean)
        else:
            signal[half_window - 1:
                   endpoint] = signal[half_window - 1:endpoint] - mean
    out = np.zeros(signal.shape)
    for i in range(window_size, signal.shape[0] - window_size):
        out[i] = fn(signal[i - half_window:i + half_window])
    return out


def running_fn_norm(signal, fn, window_size, subtract_mean=True):
    """
    Compute a normalized running function of a signal.

    Args:
        - signal (1d array): original signal
        - fn (function): the function that will be ran over windowed signal
        - window_size (int): window size, used for both running mean and
                             running function
        - subtract_mean (bool): whether the mean of the signal should be
                                subtracted before computing the running
                                function
    Returns:
        - 1d array: result of running function over the original signal
    """
    original_signal = signal.copy()
    half_window = int(window_size/2)
    endpoint = signal.shape[0] - half_window
    if subtract_mean:
        mean = running_mean_fast(signal, window_size)
        if is_odd(window_size):
            signal[half_window:endpoint] = (
                signal[half_window:endpoint] - mean)
        else:
            signal[half_window - 1:
                   endpoint] = signal[half_window - 1:endpoint] - mean
    out = np.zeros(signal.shape)
    mean = np.zeros(signal.shape)
    for i in range(window_size, signal.shape[0] - window_size):
        out[i] = fn(signal.copy()[i - half_window:i + half_window])
        mean[i] = (np.mean(original_signal[i - half_window:i + half_window] +
                   np.abs(original_signal.min())))
    return out/mean
