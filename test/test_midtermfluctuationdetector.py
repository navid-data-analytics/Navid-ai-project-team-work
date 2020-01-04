import pandas as pd
import numpy as np
from src.models import MidTermFluctuationDetector
from src.models import PeriodicityRemovalMidtermFluctuationDetector
import warnings

warnings.filterwarnings("ignore", message="Mean of empty slice")
warnings.filterwarnings("ignore", message="invalid value encountered in " +
                        "true_divide")
warnings.filterwarnings("ignore", message="invalid value encountered in " +
                        "double_scalars")
warnings.filterwarnings("ignore", message="divide by zero encountered in " +
                        "double_scalars")

DETECTORS = [
            MidTermFluctuationDetector,
            PeriodicityRemovalMidtermFluctuationDetector
            ]


def generate_input_df(value_tuples_list):
    return pd.concat([pd.DataFrame({'time': ['dummy' for i in range(val[1])],
                     'value': [val[0] for i in range(val[1])]}) for val in
                     value_tuples_list]).reset_index()


def setup_models(thresholds=(-0.1, 0.1)):
    models = []
    outs = []
    for detector in DETECTORS:
        model_out = []
        model_configuration = {1: {'thresholds': thresholds}}
        model = MidTermFluctuationDetector(model_configuration)
        model.output.connect(model_out.append)
        models.append(model)
        outs.append(model_out)
    return models, outs


def test_if_returns_none_before_memory_fill(
    extra_days=-2, expected_result={'app_id_model': 1,
                                    'time_model': 'dummy',
                                    'value_model': None}):
    models, models_outs = setup_models()
    for model, model_out in zip(models, models_outs):
        df = generate_input_df([(0, model._memory_size + extra_days)])
        for ind, row in df.iterrows():
            model.input(pd.DataFrame(row).T)
            assert model_out[ind] == expected_result


def test_pre_fill_memory_size(extra_days=-2):
    models, models_outs = setup_models()
    for model, model_out in zip(models, models_outs):
        df = generate_input_df([(1, model._memory_size + extra_days)])
        for ind, row in df.iterrows():
            model.input(pd.DataFrame(row).T)
        assert len(model_out) == len(model.memory)


def test_if_memory_doesnt_overfill(extra_days=20):
    models, models_outs = setup_models()
    for model, model_out in zip(models, models_outs):
        df = generate_input_df([(0, model._memory_size + extra_days)])
        for ind, row in df.iterrows():
            model.input(pd.DataFrame(row).T)
        assert model._memory_size == len(model.memory)


def test_if_detects_rise():
    models, models_outs = setup_models()
    for model, model_out in zip(models, models_outs):
        df = generate_input_df([(0.1, model._memory_size), (100000000, 20)])
        for ind, row in df.iterrows():
            model.input(pd.DataFrame(row).T)
        assert (model_out[model._memory_size+7]['value_model']) == 1


def test_if_detects_normal(extra_days=5, check_day=2, expected_out=0):
    models, models_outs = setup_models()
    for model, model_out in zip(models, models_outs):
        df = generate_input_df([(0.1, model._memory_size + 5)])
        for ind, row in df.iterrows():
            model.input(pd.DataFrame(row).T)
        assert ((model_out[model._memory_size +
                check_day]['value_model']) == expected_out)


def test_if_detects_decline(stable_days=28, check_day=27, expected_out=-1):
    models, models_outs = setup_models((0.5, 5))
    for model, model_out in zip(models, models_outs):
        half_memory_size = int(model._memory_size/2)
        input_values = [[(1, 1), (-1, 1)] for i in range(half_memory_size)]
        input_values = [item for sublist in input_values for item in sublist]
        input_values.append((1, stable_days))
        df = generate_input_df(input_values)
        for ind, row in df.iterrows():
            model.input(pd.DataFrame(row).T)
        assert ((model_out[model._memory_size +
                check_day]['value_model']) == expected_out)


def test_if_return_right_values(standard_expected_value=0.36340683046962396,
                                periodic_expected_value=0.06042590576789642):
    models, outs = setup_models()
    input_values = np.cos(np.linspace(0, 48, 56)) + 2
    df = pd.concat([pd.DataFrame({'time': 'dummy',
                                  'value': val}, index=[ind]) for ind, val in
                    enumerate(input_values)]).reset_index()

    for model, out in zip(models, outs):
        for ind, row in df.iterrows():
            model.input(pd.DataFrame(row).T)

    models[0]._get_fluctuation_metric_point() == standard_expected_value
    models[1]._get_fluctuation_metric_point() == periodic_expected_value
