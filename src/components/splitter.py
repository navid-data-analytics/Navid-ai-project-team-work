from src.components.pushplug import PushPlug


class Splitter:
    def __init__(self, process=lambda item: {"x": True, "y": False}):
        """
        A  splitter separates an input into multiple output.
        :param process: A process determines, which output gets
        what type of data
        """
        self._process = process
        self._outputs = {}
        self._default_output = PushPlug()
        self._touched = 0
        self._inputs_num = 0
        self._output = PushPlug()

    def get_output(self, name=None):
        """
        Get an output registered to the component with a given name
        :param name: The name of the output. The main process will use
        this name to generates the data for this output
        :return: A pushplug for this output
        """
        if name is None:
            return self._default_output
        result = self._outputs.get(name, None)
        if result is None:
            result = PushPlug()
            self._outputs.update({name: result})
        return result

    def input(self, item):
        result = self._process(item)
        for output_name, value in result.items():
            output = self.get_output(output_name)
            output.transmit(value)


class TermSplitter(Splitter):
    """Created for resending ShortTerm and MidTerm data."""

    def input(self, item):
        result = self._process(item)
        for output_name, value in result.items():
            for term in ('midterm', 'shortterm'):
                output = self.get_output('{}_{}'.format(output_name, term))
                output.transmit(value)
