package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 10,
			Up: func(db migrations.DB) error {
				logger.Info("adding RTT specific messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					INSERT INTO message_templates (type, version, template)
					VALUES
						('MidtermRttTrend15daysUp', 1, 'Your avg. Round Trip Time (RTT) per day has increased.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. RTT/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. RTT/day is {{.Number "current_score"}} ms.\nThis could be worth looking into.'),
						('MidtermRttTrend15daysDown', 1, 'Your avg. Round Trip Time (RTT) per day has decreased.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. RTT/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. RTT/day is {{.Number "current_score"}} ms.\nGreat job!'),
						('MidtermRttFluctuationImmediatelyHigh', 1, 'Your avg. Round Trip Time (RTT) per day has been fluctuating since {{.Date "current_period_start" }}. This could be worth looking into.'),
						('MidtermRttFluctuationImmediatelyStabilized', 1, 'Fluctuations in avg. Round Trip Time (RTT) per day during the period {{.Date "current_period_start"}} - {{.Date "current_period_end"}} have been successfully stabilized. Great job!');
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping RTT specific messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					DELETE FROM messages USING message_templates
						WHERE template_id IN (
							SELECT id FROM message_templates
								WHERE type IN (
									'MidtermRttTrend15daysUp',
									'MidtermRttTrend15daysDown',
									'MidtermRttFluctuationImmediatelyHigh',
									'MidtermRttFluctuationImmediatelyStabilized'
								) AND version = 1
						);
					DELETE FROM message_templates
						WHERE type IN (
							'MidtermRttTrend15daysUp',
							'MidtermRttTrend15daysDown',
							'MidtermRttFluctuationImmediatelyHigh',
							'MidtermRttFluctuationImmediatelyStabilized'
						) AND version = 1;
				`, opts.RootRole))

				return err
			},
		}
	})
}
