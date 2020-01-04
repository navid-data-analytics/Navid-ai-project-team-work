package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 13,
			Up: func(db migrations.DB) error {
				logger.Info("Re-fixing prediction period...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					BEGIN TRANSACTION;
					UPDATE message_templates SET template='Your daily average calls are expected to <span style="color:green; font-weight: bold">grow</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span> during the period {{.Date "future_period_start" }}  - {{.Date "future_period_end" }}.' WHERE type = 'MidtermPrediction15daysUp';
					UPDATE message_templates SET template='Your daily average calls are expected to <span style="color:red; font-weight: bold">decline</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span> during the period {{.Date "future_period_start" }}  - {{.Date "future_period_end" }}.' WHERE type = 'MidtermPrediction15daysDown';
					COMMIT;
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("Un-Re-fixing prediction period...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					BEGIN TRANSACTION;
					UPDATE message_templates SET template='Your daily average calls are expected to <span style="color:green; font-weight: bold">grow</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span> during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}.' WHERE type = 'MidtermPrediction15daysUp';
					UPDATE message_templates SET template='Your daily average calls are expected to <span style="color:red; font-weight: bold">decline</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span> during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}.' WHERE type = 'MidtermPrediction15daysDown';
					COMMIT;
					`, opts.RootRole))

				return err
			},
		}
	})
}
