package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 8,
			Up: func(db migrations.DB) error {
				logger.Info("adding OQ specific messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					INSERT INTO message_templates (type, version, template)
					VALUES
						('MidtermOQTrend15daysUp', 1, 'Your avg. objective quality per day has increased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nGreat job!'),
						('MidtermOQTrend15daysDown', 1, 'Your avg. objective quality per day has decreased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThis could be worth looking into.'),
						('MidtermOQFluctuationImmediatelyHigh', 1, 'Your avg. objective quality per day has been fluctuating since {{.Date "current_period_start" }}. This could be worth looking into.'),
						('MidtermOQFluctuationImmediatelyStabilized', 1, 'Fluctuations in avg. objective quality per day during the period {{.Date "current_period_start"}} - {{.Date "current_period_end"}} have been successfully stabilized. Great job!');
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping OQ specific messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					DELETE FROM messages USING message_templates
						WHERE template_id IN (
							SELECT id FROM message_templates
								WHERE type IN (
									'MidtermOQTrend15daysUp',
									'MidtermOQTrend15daysDown',
									'MidtermOQFluctuationImmediatelyHigh',
									'MidtermOQFluctuationImmediatelyStabilized'
								) AND version = 1
					  );
					DELETE FROM message_templates
						WHERE type IN (
							'MidtermOQTrend15daysUp',
							'MidtermOQTrend15daysDown',
							'MidtermOQFluctuationImmediatelyHigh',
							'MidtermOQFluctuationImmediatelyStabilized'
						) AND version = 1;
				`, opts.RootRole))

				return err
			},
		}
	})
}
