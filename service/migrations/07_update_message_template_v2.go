package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 7,
			Up: func(db migrations.DB) error {
				logger.Info("updating message templates v2...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					BEGIN TRANSACTION;
					UPDATE message_templates SET template='Your avg. no. of conferences per day has increased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. conferences/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. conferences/day.\nGreat job!' WHERE type = 'MidtermTrend15daysUp';
					UPDATE message_templates SET template='Your avg. no. of conferences per day has decreased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. conferences/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. conferences/day.\nThis could be worth looking into.' WHERE type = 'MidtermTrend15daysDown';
					UPDATE message_templates SET template='Your avg. no. of conferences per day has increased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. conferences/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. conferences/day.\nGreat job!' WHERE type = 'MidtermTrendImmediatelyUp';
					UPDATE message_templates SET template='Your avg. no. of conferences per day has decreased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. conferences/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. conferences/day.\nThis could be worth looking into.' WHERE type = 'MidtermTrendImmediatelyDown';
					UPDATE message_templates SET template='Your no. of conferences per day has been fluctuating since {{.Date "current_period_start" }}.\nThis could be worth looking into.' WHERE type = 'MidtermFluctuationImmediatelyHigh';
					UPDATE message_templates SET template='Fluctuations in the no. of conferences per day during the period {{.Date "current_period_start"}} - {{.Date "current_period_end"}} have been successfully stabilized.' WHERE type = 'MidtermFluctuationImmediatelyStabilized';
					UPDATE message_templates SET template='Your avg. no. of conferences per day is expected to grow by {{.Number "percentage"}}%% during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}).' WHERE type = 'MidtermPrediction15daysUp';
					UPDATE message_templates SET template='Your avg. no. of conferences per day is expected to decline by {{.Number "percentage"}}%% during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}.' WHERE type = 'MidtermPrediction15daysDown';
					COMMIT;
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping message templates v2...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
				`, opts.RootRole))

				return err
			},
		}
	})
}
