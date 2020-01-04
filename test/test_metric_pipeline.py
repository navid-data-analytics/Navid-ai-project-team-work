from src.pipeline import MetricPipeline
from src.pipeline import AppPipeline
from src.components import Transceiver

app_ids = [0, 1, 2]
input_dict = {0: 'a',
              1: 'a',
              2: 'a'}


# NOTE: The reason for mocking app pipeline is that MetricPipeline
#       initializes it inside itself.
class MockAppPipeline(AppPipeline):
    def __init__(self, config, appID, AidServiceConnection=None):
        super(MockAppPipeline, self).__init__(config, appID,
                                              AidServiceConnection)

    def _joiner_process(self, values):
        return {'a': values['a'],
                'b': values['b'],
                'c': values['c'], }

    def _link_components(self):
        self._entrypoint.output.connect(self._a_model.input)
        self._a_model.output.connect(self._a_decision.input)
        self._a_decision.output.connect(self._exitpoint.get_input('a'))

        self._entrypoint.output.connect(self._b_model.input)
        self._b_model.output.connect(self._b_decision.input)
        self._b_decision.output.connect(self._exitpoint.get_input('b'))

        self._entrypoint.output.connect(self._c_model.input)
        self._c_model.output.connect(self._c_decision.input)
        self._c_decision.output.connect(self._exitpoint.get_input('c'))

    def _create_components(self, config=None):
        self._a_model = Transceiver()
        self._b_model = Transceiver()
        self._c_model = Transceiver()
        self._a_decision = Transceiver()
        self._b_decision = Transceiver()
        self._c_decision = Transceiver()


def test_metric_pipeline_builds():
    MetricPipeline(MockAppPipeline, app_ids, config=None)


def test_metric_pipeline_runs():
    result = []
    pipeline = MetricPipeline(MockAppPipeline, app_ids, config=None)
    pipeline.output.connect(result.append)
    pipeline.input(input_dict)
    assert list(result[0].keys()) == app_ids
