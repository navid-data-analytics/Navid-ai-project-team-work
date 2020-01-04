package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 9,
			Up: func(db migrations.DB) error {
				logger.Info("updating prediction message templates v3...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					BEGIN TRANSACTION;
					UPDATE message_templates SET template='Your avg. no. of conferences per day is expected to grow by {{.Number "percentage"}}%% during the period {{.Date "future_period_start" }}  - {{.Date "future_period_end" }}).' WHERE type = 'MidtermPrediction15daysUp';
					UPDATE message_templates SET template='Your avg. no. of conferences per day is expected to decline by {{.Number "percentage"}}%% during the period {{.Date "future_period_start" }}  - {{.Date "future_period_end" }}.' WHERE type = 'MidtermPrediction15daysDown';
					COMMIT;
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping prediction message templates v3...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					BEGIN TRANSACTION;
					UPDATE message_templates SET template='Your avg. no. of conferences per day is expected to grow by {{.Number "percentage"}}%% during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}).' WHERE type = 'MidtermPrediction15daysUp';
					UPDATE message_templates SET template='Your avg. no. of conferences per day is expected to decline by {{.Number "percentage"}}%% during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}.' WHERE type = 'MidtermPrediction15daysDown';
					COMMIT;
				`, opts.RootRole))

				return err
			},
		}
	})
}
