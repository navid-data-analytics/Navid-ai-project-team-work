from src.components.pushplug import PushPlug


class Broadcaster:
    def __init__(self, process=lambda item: item):
        """
        Receives an input and forwards it on every output it has.
        :param process: performs a transform operation if it is necessary
        """
        self._process = process
        self._outputs = []

    @property
    def output(self):
        result = PushPlug()
        self._outputs.append(result)
        return result

    def input(self, value):
        for output in self._outputs:
            item = self._process(value)
            output.transmit(item)
