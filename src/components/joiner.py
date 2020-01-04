from src.components.pushplug import PushPlug


class Joiner:
    def __init__(self, process=lambda inputs: True):
        """
        Collects data on it's inputs and joins them.
        Note: The join operation triggered only if all input received at least
        one data
        :param process: The joiner process receives the inputs as dictionary
        and returns an output transmitted towards
        """
        self._process = process
        self._inputs = {}
        self._touches = {}
        self._inputs_num = 0
        self._output = PushPlug()

    @property
    def output(self):
        return self._output

    def _check(self):
        for name, touches in self._touches.items():
            if touches < 1:
                return False
        return True

    def _insert(self, name, value):
        items = self._inputs.get(name, [])
        touches = self._touches.get(name, 0)
        items.append(value)
        touches = self._touches.get(name, 0)
        self._touches.update({name: touches + 1})
        self._inputs.update({name: items})

        for input_name in self._touches.keys():
            touches = self._touches.get(input_name, 0)
            if touches < 1:
                return
        inputs = {}
        forward_none = True
        for input_name in self._inputs.keys():
            items = self._inputs.get(input_name)
            touches = self._touches.get(input_name, 0)
            input_ = items.pop(0)
            if input_ is not None:
                forward_none = False
            inputs.update({input_name: input_})
            self._inputs.update({input_name: items})
            self._touches.update({input_name: touches - 1})

        if forward_none is True:
            self._output.transmit(None)
            return
        output = self._process(inputs)
        self._output.transmit(output)

    def get_input(self, name):
        self._touches[name] = 0
        return lambda value: self._insert(name, value)
