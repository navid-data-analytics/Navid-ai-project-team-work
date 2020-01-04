from src.decisionmakers import DecisionMaker
from src.utils import measure_time_metric
from prometheus_client import Gauge
from collections import defaultdict
import logging

trend_direction_gauge = Gauge('complex_DM_direction', 'metric',
                              ['app_id', 'metric_name'])
trend_type_gauge = Gauge('complex_DM_trend_type', 'metric',
                         ['app_id', 'metric_name'])
trend_reasons_gauge = Gauge('complex_DM_reasons', 'metric',
                            ['app_id', 'metric_name'])

logger = logging.getLogger('root')


class ComplexOQDecisionMaker(DecisionMaker):
    def __init__(self,
                 app_ids,
                 metric='objective_quality_v35_average',
                 support_metrics=['lossEffectMean', 'throughputEffectMean',
                                  'delayEffectMean'],
                 aid_service_connection=None):
        super(ComplexOQDecisionMaker, self).__init__(
            app_id=app_ids,
            metric=metric,
            aid_service_connection=aid_service_connection,
            process=self._run)
        self._support_metrics = support_metrics
        self._reset_decision_flags()
        self._immediate_reasons = {
            app_id: {
                1: [],
                -1: []
            }
            for app_id in app_ids
        }
        self._reasons = defaultdict(list)

    @measure_time_metric
    def _run(self, dms_output):
        """
        Run the main function for the complex dm component.

        Arguments:
        - dms_output: a dictionary of a form:
        {'MAIN METRIC': {},
         'SUPPORT METRIC1': {},
         'SUPPORT METRIC2': {},
         ...
         },
         where each value for metrics is the following dict:
            'time_model' - Time,
            'value_model'- verdict
            'app_id_model' - appID
            'message_details' - message arguments for gRPC
            'trend_monthly_decision' - flag for monthly trend
            'trend_immediate_decision' - flag for immediate trend

        Returns:
        - output: dict with the following fields:
            'message_details' - message arguments for gRPC,
            'type': type of trend (immediate, monthly),
            'reasons': list of reasons for change in OQ,
            'direction': upwards or downwards
        """

        logger.debug(repr(self) + 'Starting _run method.')
        logger.debug(
            repr(self) +
            'Input of Complex Decision Maker: {}'.format(dms_output))
        output = {}
        for term in ('shortterm', 'midterm'):
            preprocessed_output = self._preprocess(dms_output, term)
            date = self._get_date(preprocessed_output)
            decision = self._decide(preprocessed_output)
            self._postprocess(decision)
            output[term] = {'complex_dm_output': decision, 'date': date}

        logger.debug(
            repr(self) + 'Complex Decision Maker Output: {}.'.format(output))
        return output

    def _preprocess(self, dms_output, term):
        """ Unpack the input and prepare it to decide"""
        self._handle_missing_keys(dms_output, term)
        unpacked_output = self._unpack_output(dms_output, term)
        return unpacked_output

    def _handle_missing_keys(self, dms_output, term):
        for metric in dms_output.keys():
            if metric in [
                    'rtt_average_{}'.format(term),
                    'conferences_terminated_{}'.format(term)
            ]:
                continue
            for app_id in self.app_ids:
                if 'message_details' not in list(
                        dms_output[metric][app_id]['trend_detection'].keys()):
                    dms_output[metric][app_id]['trend_detection'][
                        'message_details'] = (None, )
                if 'sent' not in list(
                        dms_output[metric][app_id]['trend_detection'].keys()):
                    dms_output[metric][app_id]['trend_detection'][
                        'sent'] = False

    def _check_main_type_and_direction(self, main_metric_decision):
        if 'trend_monthly_decision' not in list(
                main_metric_decision['trend_detection'].keys()):
            return
        self._determine_trend_type(main_metric_decision)
        self._determine_trend_direction(main_metric_decision)

    def _determine_trend_type(self, main_metric_decision):
        logger.debug('Setting trend type')
        if main_metric_decision['trend_detection']['trend_monthly_decision']:
            logger.debug('Setting trend type to monthly!')
            self._trend_type = 'monthly'
        elif main_metric_decision['trend_detection'][
                'trend_immediate_decision']:
            logger.debug('Setting trend type to immediate!')
            self._trend_type = 'immediate'
        logger.debug('Done setting trend type to {}'.format(self._trend_type))

    def _determine_trend_direction(self, main_metric_decision):
        logger.debug('Determining trend direction!')
        if not self._trend_type:
            logger.debug('Cannot determine trend direction!')
            return
        self._direction = main_metric_decision['trend_detection'][
            'value_model']
        logger.debug('Direction set as {}'.format(self._direction))

    def _unpack_output(self, dms_output, term):
        main_metric, support_metrics = self._split_metrics(dms_output, term)
        return main_metric, support_metrics

    def _split_metrics(self, outputs, term):
        main_metric = outputs[self.metric + '_{}'.format(term)]
        del outputs[self.metric + '_{}'.format(term)]
        support_metrics = self._get_cleaned_support_metrics(outputs.copy(),
                                                            term)
        return main_metric, support_metrics

    def _get_date(self, input):
        app_ids = input[0].keys()
        for app_id in app_ids:
            for operation in input[0][app_id].keys():
                date = input[0][app_id][operation].get('time_model', None)
                if date is not None:
                    return date

    def _decide(self, preprocessed_output):
        main_metric, support_metrics = preprocessed_output
        output = defaultdict(dict)
        self._handle_decisions(main_metric, support_metrics, output)
        return output

    def _handle_decisions(self, main_metric, support_metrics, output):
        for app_id in self.app_ids:
            self._handle_decision_for_app(main_metric, support_metrics, app_id,
                                          output)
            self._handle_prometheus(app_id)
            self._reset_decision_flags()

    def _handle_decision_for_app(self, main_metric, support_metrics, app_id,
                                 output):
        self._check_main_type_and_direction(main_metric[app_id])
        self._find_immediate_reasons(support_metrics, app_id)
        if self._trend_type and self._direction:
            self._find_complex_reasons(support_metrics, app_id)
        decision = self._prepare_decision(main_metric[app_id], app_id)
        output[self.metric][app_id] = {'trend_detection': decision}

    def _get_cleaned_support_metrics(self, outputs, term):
        valid_support_metrics = ['{}_{}'.format(metric, term) for metric in
                                 self._support_metrics]
        support_metrics = {}
        for name in valid_support_metrics:
            metric = outputs.get(name, None)
            if metric is not None:
                support_metrics[name] = outputs[name]
        return support_metrics

    def _prepare_decision(self, main_metric, appid):
        self._check_if_sendable(appid)
        output = {
            'message_details':
            main_metric['trend_detection']['message_details'],
            'type': self._trend_type,
            'reasons': sorted(self._reasons[appid].copy()),
            'send': self._send,
            'direction': self._direction
        }
        return output

    def _check_if_sendable(self, appid):
        correct_type = self._trend_type in ['monthly', 'immediate']
        correct_type &= self._direction != 0
        correct_reason = len(self._reasons[appid]) > 0
        self._check_OQ_reasons_mismatch(correct_type, correct_reason, appid)
        self._send = correct_type and correct_reason

    def _check_OQ_reasons_mismatch(self, correct_type, correct_reason, appid):
        if correct_type and not correct_reason:
            if len(self._immediate_reasons[appid][self._direction]) == 0:
                logger.error(
                    "The OQ changed, but reasons were not detected!\ndirection: {}\ntrend type: {}"  # noqa
                    .format(self._direction, self._trend_type))
        elif correct_type and correct_reason:
            self._reasons[appid].insert(0, self.metric)

    def _find_complex_reasons(self, support_metrics, appid):
        if self._trend_type != 'monthly':
            self._reasons[appid] = list(
                set(self._immediate_reasons[appid][self._direction]))
            return
        self._search_reasons_for_monthly(support_metrics, appid)

    def _search_reasons_for_monthly(self, support_metrics, appid):
        for metric in support_metrics.keys():
            if metric in ['rtt_average', 'conferences_terminated']:
                continue
            support_direction = support_metrics[metric][appid][
                'trend_detection']['value_model']
            self._append_reason(support_direction, metric, appid)

    def _append_reason(self, support_direction, metric, appid):
        should_append = self._check_append_conditions(support_direction)
        if should_append:
            logger.debug(repr(self) + 'appending reason {}.'.format(metric))
            self._reasons[appid].append(metric)

    def _check_append_conditions(self, support_direction):
        valid_verdict = self._direction != 0
        matching_directions = support_direction == self._direction
        should_append = valid_verdict and matching_directions
        return should_append

    def _find_immediate_reasons(self, support_metrics, appid):
        for metric in support_metrics.keys():
            if metric in ['rtt_average', 'conferences_terminated']:
                continue
            if self._trend_type == 'immediate':
                self._append_immediate_reason(support_metrics, metric, appid)

    def _append_immediate_reason(self, support_metrics, metric, appid):
        direction = support_metrics[metric][appid]['trend_detection'][
            'value_model']
        if direction != 0:
            logger.debug(
                repr(self) +
                'im reason {} of value {}.'.format(metric, direction))
            self._immediate_reasons[appid][direction].append(metric)

    def _postprocess(self, output):
        """Flush the state after processing data."""
        self._flush_immediate_reasons(output)
        self._flush_monthly_reasons()
        self._reset_decision_flags()

    def _handle_prometheus(self, app_id):
        trend_type = 1 if self._trend_type == 'immediate' else 2
        direction = 0 if self._direction is None else self._direction
        trend_direction_gauge.labels(app_id, self.metric).set(direction)
        trend_type_gauge.labels(app_id, self.metric).set(trend_type)
        trend_reasons_gauge.labels(app_id, self.metric).set(
            len(self._reasons[app_id]))

    def _flush_immediate_reasons(self, output):
        for app_id in self.app_ids:
            if output[self.metric][app_id]['trend_detection'][
                    'type'] != 'monthly':  # noqa
                continue
            self._flush_immediate_reasons_per_app()

    def _flush_immediate_reasons_per_app(self):
        for app_id in self.app_ids:
            self._immediate_reasons[app_id][1].clear()
            self._immediate_reasons[app_id][-1].clear()

    def _flush_monthly_reasons(self):
        for app_id in self.app_ids:
            self._reasons[app_id].clear()

    def _reset_decision_flags(self):
        self._trend_type = None
        self._direction = None
        self._send = False

    @property
    def app_ids(self):
        return self.app_id

    def __repr__(self):
        """Use repr(self) for DEBUG, ERROR, CRITICAL logging levels."""
        return 'Metric {}: '.format(self.metric)

    def __str__(self):
        """Use str(self) for INFO, WARNING logging levels."""
        return 'Metric {}: '.format(self.metric)
