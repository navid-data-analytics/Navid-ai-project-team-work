from src.components.pushplug import PushPlug


class Transmitter:
    def __init__(self, process=lambda: True):
        """
        A source component generates and transmits data
        :param process: the producer process generates the desired data
        """
        self._process = process
        self._output = PushPlug()

    @property
    def output(self):
        return self._output

    def run(self):
        """
        Trigger an execution of the Transmitter process and
        transmits the generated data to the output.
        """
        item = self._process()
        self._output.transmit(item)
