package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 16,
			Up: func(db migrations.DB) error {
				logger.Info("adding immediate short messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					INSERT INTO message_templates (type, version, template)
					VALUES
						('ShorttermTrendImmediatelyUp', 1, 'Your daily average calls have <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\n\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} average calls/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} average calls/day.\n\nGreat job!'),
						('ShorttermTrendImmediatelyDown', 1, 'Your daily average calls have <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\n\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} average calls/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} average calls/day.\n\nThis could be worth looking into.'),
						('ShorttermRttTrendImmediatelyUp', 1, 'Your avg. Round Trip Time (RTT) per day has <span style="color:red; font-weight: bold">increased</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. RTT/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. RTT/day is {{.Number "current_score"}} ms.\nThis could be worth looking into.'),
						('ShorttermRttTrendImmediatelyDown', 1, 'Your avg. Round Trip Time (RTT) per day has <span style="color:green; font-weight: bold">decreased</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. RTT/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. RTT/day is {{.Number "current_score"}} ms.\nGreat job!'),
						('ShorttermRttFluctuationImmediatelyHigh', 1, 'Your avg. Round Trip Time (RTT) per day has been fluctuating since {{.Date "current_period_start" }}.\nThis could be worth looking into.'),
						('ShorttermRttFluctuationImmediatelyStabilized', 1, 'Fluctuations in the avg. Round Trip Time (RTT) per day during the period {{.Date "current_period_start"}} - {{.Date "current_period_end"}} have been successfully stabilized.'),
						('ShorttermOQFluctuationImmediatelyHigh', 1, 'Your avg. objective quality per day has been fluctuating since {{.Date "current_period_start" }}.\nThis could be worth looking into.'),
						('ShorttermOQFluctuationImmediatelyStabilized', 1, 'Fluctuations in the avg. objective quality per day during the period {{.Date "current_period_start"}} - {{.Date "current_period_end"}} have been successfully stabilized.'),
						('ShorttermFluctuationImmediatelyHigh', 1, 'Your no. of conferences per day has been fluctuating since {{.Date "current_period_start" }}.\nThis could be worth looking into.'),
						('ShorttermFluctuationImmediatelyStabilized', 1, 'Fluctuations in the no. of conferences per day during the period {{.Date "current_period_start"}} - {{.Date "current_period_end"}} have been successfully stabilized.'),
						('ShorttermPrediction7daysDown', 1, 'Your daily average calls are expected to <span style="color:red; font-weight: bold">decline</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span> during the period {{.Date "future_period_start" }}  - {{.Date "future_period_end" }}.'),
						('ShorttermPrediction7daysUp', 1, 'Your daily average calls are expected to <span style="color:green; font-weight: bold">grow</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span> during the period {{.Date "future_period_start" }}  - {{.Date "future_period_end" }}.');
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping immediate short messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					DELETE FROM messages USING message_templates
						WHERE template_id IN (
							SELECT id FROM message_templates
								WHERE type IN (
									'ShorttermTrendImmediatelyUp',
									'ShorttermTrendImmediatelyDown',
									'ShorttermRttTrendImmediatelyDown',
									'ShorttermRttTrendImmediatelyUp',
									'ShorttermRttFluctuationImmediatelyHigh',
									'ShorttermRttFluctuationImmediatelyStabilized',
									'ShorttermOQFluctuationImmediatelyHigh',
									'ShorttermOQFluctuationImmediatelyStabilized',
									'ShorttermFluctuationImmediatelyHigh',
									'ShorttermFluctuationImmediatelyStabilized',
									'ShorttermPrediction7daysDown',
									'ShorttermPrediction7daysUp'
								) AND version = 1
						);
					DELETE FROM message_templates
						WHERE type IN (
							'ShorttermTrendImmediatelyUp',
							'ShorttermTrendImmediatelyDown',
							'ShorttermRttTrendImmediatelyDown',
							'ShorttermRttTrendImmediatelyUp',
							'ShorttermRttFluctuationImmediatelyHigh',
							'ShorttermRttFluctuationImmediatelyStabilized',
							'ShorttermOQFluctuationImmediatelyHigh',
							'ShorttermOQFluctuationImmediatelyStabilized',
							'ShorttermFluctuationImmediatelyHigh',
							'ShorttermFluctuationImmediatelyStabilized',
							'ShorttermPrediction7daysDown',
							'ShorttermPrediction7daysUp'
						) AND version = 1;
				`, opts.RootRole))

				return err
			},
		}
	})
}
