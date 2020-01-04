package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 5,
			Up: func(db migrations.DB) error {
				logger.Info("adding message templates v1...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					INSERT INTO message_templates (type, version, template)
					VALUES
						('MidtermTrend15daysUp', 1, 'Your daily average calls have increased by {{.Number "percentage"}}%% in the last 30 days. Great job!'),
						('MidtermTrend15daysDown', 1, 'Your daily average calls have decreased by {{.Number "percentage"}}%% in the last 30 days. This could be worth looking into.'),
						('MidtermTrendImmediatelyUp', 1, 'Your daily average calls have been growing at a rate of {{.Number "percentage"}}%% for the past {{.Number "days"}} days. Great job!'),
						('MidtermTrendImmediatelyDown', 1, 'Your daily average calls is decreasing at a rate of {{.Number "percentage"}}%% for the past {{.Number "days"}} days. This could be worth looking into.'),
						('MidtermFluctuationImmediatelyHigh', 1, 'Your service usage has been fluctuating in the past 30 days. This could be worth looking into.'),
						('MidtermFluctuationImmediatelyStabilized', 1, 'Your service usage fluctuations have been successfully stabilized. Great job!'),
						('MidtermPrediction15daysUp', 1, 'Your daily average calls is expected to grow by {{.Number "percentage"}}%% in the next 30 days. Great job!'),
						('MidtermPrediction15daysDown', 1, 'Your daily average calls is expected to decline by {{.Number "percentage"}}%% in the next 30 days. You might want to check this out.');
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping message templates v1...")
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
								) AND version = 1
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
						) AND version = 1;
				`, opts.RootRole))

				return err
			},
		}
	})
}
