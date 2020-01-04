package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 6,
			Up: func(db migrations.DB) error {
				logger.Info("adding message templates v2...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					INSERT INTO message_templates (type, version, template)
					VALUES
						('MidtermTrend15daysUp', 2, 'Your daily average conferences have increased by {{.Number "percentage"}}%%\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. conferences/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. conferences/day.\nGreat job!'),
						('MidtermTrend15daysDown', 2, 'Your daily average calls have decreased by {{.Number "percentage"}}%%\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. conferences/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. conferences/day.\nThis could be worth looking into.'),						
						('MidtermTrendImmediatelyUp', 2, 'Your daily average conferences have increased by {{.Number "percentage"}}%%\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. conferences/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. conferences/day.\nGreat job!'),
						('MidtermTrendImmediatelyDown', 2, 'Your daily average calls have decreased by {{.Number "percentage"}}%%\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. conferences/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. conferences/day.\nThis could be worth looking into.'),
						('MidtermFluctuationImmediatelyHigh', 2, 'Your daily average conference count has been fluctuating since {{.Date "current_period_start" }}.\nThis could be worth looking into.'),
						('MidtermFluctuationImmediatelyStabilized', 2, 'Your daily average calls count fluctuations during the period {{.Date "current_period_start"}} - {{.Date "current_period_end"}} have been successfully stabilized.\nGreat job!'),
						('MidtermPrediction15daysUp', 2, 'Your daily average conferences are expected to grow by {{.Number "percentage"}}%% during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}).\nGood luck!'),
						('MidtermPrediction15daysDown', 2, 'Your daily average calls are expected to decline by {{.Number "percentage"}}%% during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}.\nYou might want to keep an eye on it.');
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping message templates v2...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					DELETE FROM messages USING message_templates
						WHERE template_id IN (
							SELECT id FROM message_templates
								WHERE type IN (
									'MidtermTrend15daysUp',
									'MidtermTrend15daysDown',
									'MidtermTrendImmediatelyUp',
									'MidtermTrendImmediatelyDown',
									'MidtermFluctuationImmediatelyHigh',
									'MidtermFluctuationImmediatelyStabilized',
									'MidtermPrediction15daysUp',
									'MidtermPrediction15daysDown'
								) AND version = 2
					  );
					DELETE FROM message_templates
						WHERE type IN (
							'MidtermTrend15daysUp',
							'MidtermTrend15daysDown',
							'MidtermTrendImmediatelyUp',
							'MidtermTrendImmediatelyDown',
							'MidtermFluctuationImmediatelyHigh',
							'MidtermFluctuationImmediatelyStabilized',
							'MidtermPrediction15daysUp',
							'MidtermPrediction15daysDown'
						) AND version = 2;
				`, opts.RootRole))

				return err
			},
		}
	})
}
