from src.components import Transceiver


class Model(Transceiver):
    """ An abstract class of the model component. """

    def __init__(self, app_id=1, process=lambda: None, metric=None):
        """
        Construct a Model object.

        Arguments:
        - app_id: int, self-explainatory
        - process: function run upon data arrival
        """
        self._app_id = app_id
        self._metric = metric
        Transceiver.__init__(self, process=process)
        self.filled = lambda x: len(x) == x.maxlen

    def _run(self, input_signal):
        raise NotImplementedError('Not Implemented!')

    def _predict(self):
        raise NotImplementedError('Not Implemented!')

    def _postprocess(self):
        raise NotImplementedError('Not Implemented!')

    @property
    def app_id(self):
        return self._app_id

    @property
    def metric(self):
        return self._metric
