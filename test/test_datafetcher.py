import datetime
from src.components import DatePackage
from src.pipeline import DataFetcher
import pandas as pd


UTC = datetime.timezone.utc
DAY_WITH_DATA = datetime.datetime(2019, 2, 12, tzinfo=UTC)
# tests below are only checking backfill loading functionality, which does not
# utilize end date, hence end date is set to "ANY_DAY" value
ANY_DAY = datetime.datetime(2004, 9, 11, tzinfo=UTC)
METRICS = ['delayEffectMean', 'throughputEffectMean', 'lossEffectMean']
APPIDS = [331463193,
          347489791,
          380077084,
          486784909,
          722018081,
          748139165,
          815616092,
          739234857,
          928873129,
          943549171,
          234913325,
          602363212]

DAY_WITH_DATA_VALUES = {
    'delayEffectMean': {
        331463193:
        {'value': 0.97258612869268202, 'time': pd.to_datetime('2019-02-12')},
        347489791:
        {'value': 0.91537794037105613, 'time': pd.to_datetime('2019-02-12')},
        380077084:
        {'value': 0.98760898014324461, 'time': pd.to_datetime('2019-02-12')},
        486784909:
        {'value':  0.98701524312203148, 'time': pd.to_datetime('2019-02-12')},
        722018081:
        {'value': 0.98276524392430764, 'time': pd.to_datetime('2019-02-12')},
        748139165:
        {'value': 0, 'time': pd.to_datetime('2019-02-12')},
        815616092:
        {'value': 0.92702925484376009, 'time': pd.to_datetime('2019-02-12')},
        739234857:
        {'value': 0.95307033580861367, 'time': pd.to_datetime('2019-02-12')},
        928873129:
        {'value': 0.98291775273266802, 'time': pd.to_datetime('2019-02-12')},
        943549171:
        {'value': 0, 'time': pd.to_datetime('2019-02-12')},
        234913325:
        {'value': 0.91630050591002143, 'time': pd.to_datetime('2019-02-12')},
        602363212:
        {'value': 0.9757279250780718, 'time': pd.to_datetime('2019-02-12')}},
    'throughputEffectMean': {
        331463193:
        {'value': 1.3980109028280003, 'time': pd.to_datetime('2019-02-12')},
        347489791:
        {'value': 1.7379735890275234, 'time': pd.to_datetime('2019-02-12')},
        380077084:
        {'value': 1.4059120631855446, 'time': pd.to_datetime('2019-02-12')},
        486784909:
        {'value': 1.6379179368811776, 'time': pd.to_datetime('2019-02-12')},
        722018081:
        {'value': 1.6957275614495899, 'time': pd.to_datetime('2019-02-12')},
        748139165:
        {'value': 0, 'time': pd.to_datetime('2019-02-12')},
        815616092:
        {'value': 1.5686380317617186, 'time': pd.to_datetime('2019-02-12')},
        739234857:
        {'value': 0.88455822317957178, 'time': pd.to_datetime('2019-02-12')},
        928873129:
        {'value': 1.4188631237147959, 'time': pd.to_datetime('2019-02-12')},
        943549171:
        {'value': 0, 'time': pd.to_datetime('2019-02-12')},
        234913325:
        {'value': 1.6458263782260343, 'time': pd.to_datetime('2019-02-12')},
        602363212:
        {'value': 1.6376799654796752, 'time': pd.to_datetime('2019-02-12')}},
    'lossEffectMean': {
        331463193:
        {'value': 0.99536056540312234, 'time': pd.to_datetime('2019-02-12')},
        347489791:
        {'value': 0.99960761445132584, 'time': pd.to_datetime('2019-02-12')},
        380077084:
        {'value': 0.99887517539247994, 'time': pd.to_datetime('2019-02-12')},
        486784909:
        {'value': 0.99999933460777812, 'time': pd.to_datetime('2019-02-12')},
        722018081:
        {'value': 1.0, 'time': pd.to_datetime('2019-02-12')},
        748139165:
        {'value': 0, 'time': pd.to_datetime('2019-02-12')},
        815616092:
        {'value': 0.98906753925816437, 'time': pd.to_datetime('2019-02-12')},
        739234857:
        {'value': 1.0, 'time': pd.to_datetime('2019-02-12')},
        928873129:
        {'value': 0.99995173234259382, 'time': pd.to_datetime('2019-02-12')},
        943549171:
        {'value': 0, 'time': pd.to_datetime('2019-02-12')},
        234913325:
        {'value': 1.0, 'time': pd.to_datetime('2019-02-12')},
        602363212:
        {'value': 0.995108984107196, 'time': pd.to_datetime('2019-02-12')}}
}


class MockCRS:
    def Aggregate(*args, **kwargs):
        return {}


def get_datafetcher_output(datepackage):
    crs = MockCRS()
    df = DataFetcher([], [], crs)
    return df._receive(datepackage)


def test_if_effects_data_fetched_on_present_day():
    datepackage = DatePackage(DAY_WITH_DATA, ANY_DAY)
    df_output = get_datafetcher_output(datepackage)
    assert list(df_output.keys()) == METRICS
    assert list(df_output[METRICS[0]]) == APPIDS

    for metric in METRICS:
        for appid in APPIDS:
            expected_day_one_value = \
                DAY_WITH_DATA_VALUES[metric][appid]['value']
            actual_day_one_value = df_output[metric][appid]['value'].values[0]
            assert actual_day_one_value == expected_day_one_value

            expected_day_one_date =  \
                DAY_WITH_DATA_VALUES[metric][appid]['time']
            actual_day_one_date = pd.to_datetime(
                                    df_output[metric][appid]['time'].values[0])
            assert actual_day_one_date == expected_day_one_date


def test_if_effects_data_not_fetched_on_early_day():
    datepackage = DatePackage(ANY_DAY - datetime.timedelta(days=1),
                              ANY_DAY)
    assert get_datafetcher_output(datepackage) == {}
