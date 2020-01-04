from src.components.pushplug import PushPlug


class Transceiver():
    def __init__(self, process=lambda value: value):
        """
        A transceiver is comprising both a transmitter and a receiver
        functionality
        :param process: A process transforms the input data to the
        desired output
        """
        self._process = process
        self._output = PushPlug()

    @property
    def output(self):
        return self._output

    def _get_process(self):
        return self._process

    def input(self, value):
        item = self._process(value)
        self._output.transmit(item)
