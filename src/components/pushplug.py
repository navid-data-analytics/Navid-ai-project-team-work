class PushPlug:
    def __init__(self):
        self._target = None

    def connect(self, target):
        self._target = target

    def is_connected(self):
        return self._target is not None

    def transmit(self, value):
        target = self._target
        if target is None:
            return
        target(value)
