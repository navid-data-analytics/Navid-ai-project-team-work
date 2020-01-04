from src.components.pushplug import PushPlug


class Collector:
    """Collect the stream of items. Upon getting 'None'
       the collector transits output further."""

    def __init__(self):
        """
        Collect items receives on it's input and forwards them as a list
        once it receives a None.
        """
        self._items = []
        self._output = PushPlug()

    @property
    def output(self):
        return self._output

    def input(self, value):
        if value is None:
            self._output.transmit(self._items)
            return
        self._items.append(value)
