class Receiver:
    def __init__(self, process=lambda value: value):
        """
        Receives a data on it's input and process it
        :param process: processor triggered to run every time an input is fed
        """
        self._process = process

    def input(self, value):
        self._process(value)
