package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 12,
			Up: func(db migrations.DB) error {
				logger.Info("Add color and bold to parts of messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					BEGIN TRANSACTION;
					UPDATE message_templates SET template='Your daily average calls have <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\n\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} average calls/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} average calls/day.\n\nGreat job!' WHERE type = 'MidtermTrend15daysUp';
					UPDATE message_templates SET template='Your daily average calls have <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\n\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} average calls/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} average calls/day.\n\nThis could be worth looking into.' WHERE type = 'MidtermTrend15daysDown';
					UPDATE message_templates SET template='Your daily average calls have <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\n\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} average calls/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} average calls/day.\n\nGreat job!' WHERE type = 'MidtermTrendImmediatelyUp';
					UPDATE message_templates SET template='Your daily average calls have <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\n\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} average calls/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} average calls/day.\n\nThis could be worth looking into.' WHERE type = 'MidtermTrendImmediatelyDown';
					UPDATE message_templates SET template='Your daily average calls count has been <span style="color:red; font-weight: bold">fluctuating</span> since <span style="font-weight: bold">{{.Date "current_period_start" }}</span>.\n\nThis could be worth looking into.' WHERE type = 'MidtermFluctuationImmediatelyHigh';
					UPDATE message_templates SET template='Your daily average calls count fluctuations during the period <span style="font-weight: bold">{{.Date "current_period_start"}} - {{.Date "current_period_end"}}</span> have been successfully <span style="color:green; font-weight: bold">stabilized</span>.' WHERE type = 'MidtermFluctuationImmediatelyStabilized';
					UPDATE message_templates SET template='Your daily average calls are expected to <span style="color:green; font-weight: bold">grow</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span> during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}.' WHERE type = 'MidtermPrediction15daysUp';
					UPDATE message_templates SET template='Your daily average calls are expected to <span style="color:red; font-weight: bold">decline</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span> during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}.' WHERE type = 'MidtermPrediction15daysDown';

					UPDATE message_templates SET template='Your average objective quality per day has <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\n\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) average OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) average OQ/day is {{.Number "current_score"}}.\n\nGreat job!' WHERE type = 'MidtermOQTrend15daysUp';
					UPDATE message_templates SET template='Your average objective quality per day has <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\n\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) average OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) average OQ/day is {{.Number "current_score"}}.\n\nThis could be worth looking into.' WHERE type = 'MidtermOQTrend15daysDown';
					UPDATE message_templates SET template='Your average objective quality per day has been <span style="color:red; font-weight: bold">fluctuating</span> since <span style="font-weight: bold">{{.Date "current_period_start" }}</span>.\n\nThis could be worth looking into.' WHERE TYPE = 'MidtermOQFluctuationImmediatelyHigh';
					UPDATE message_templates SET template='Fluctuations in average objective quality per day during the period <span style="font-weight: bold">{{.Date "current_period_start"}} - {{.Date "current_period_end"}}</span> have been successfully <span style="color:green; font-weight: bold">stabilized</span>.\n\nGreat job!' WHERE type = 'MidtermOQFluctuationImmediatelyStabilized';

					UPDATE message_templates SET template='Your average round-trip time (RTT) per day has  <span style="color:red; font-weight: bold">increased</span>.\n\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) average RTT/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. RTT/day is {{.Number "current_score"}} ms.\n\nThis could be worth looking into.' WHERE type = 'MidtermRttTrend15daysUp';
					UPDATE message_templates SET template='Your average round-trip time (RTT) per day has <span style="color:green; font-weight: bold">decreased</span>.\n\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) average RTT/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. RTT/day is {{.Number "current_score"}} ms.\n\nGreat job!' WHERE type = 'MidtermRttTrend15daysDown';
					UPDATE message_templates SET template='Your average round-trip time (RTT) per day has been <span style="color:red; font-weight: bold">fluctuating</span> since <span style="font-weight: bold">{{.Date "current_period_start" }}</span>.\n\nThis could be worth looking into.' WHERE type = 'MidtermRttFluctuationImmediatelyHigh';
					UPDATE message_templates SET template='Fluctuations in average round-trip time (RTT) per day during the period <span style="font-weight: bold">{{.Date "current_period_start"}} - {{.Date "current_period_end"}}</span> have been successfully <span style="color:green; font-weight: bold">stabilized</span>.\n\nGreat job!' WHERE type = 'MidtermRttFluctuationImmediatelyStabilized';

					COMMIT;
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("Remove color and bold from messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					BEGIN TRANSACTION;
					UPDATE message_templates SET template='Your avg. no. of calls per day has increased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. calls/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. calls/day.\nGreat job!' WHERE type = 'MidtermTrend15daysUp';
					UPDATE message_templates SET template='Your avg. no. of calls per day has decreased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. calls/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. calls/day.\nThis could be worth looking into.' WHERE type = 'MidtermTrend15daysDown';
					UPDATE message_templates SET template='Your avg. no. of calls per day has increased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. calls/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. calls/day.\nGreat job!' WHERE type = 'MidtermTrendImmediatelyUp';
					UPDATE message_templates SET template='Your avg. no. of calls per day has decreased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }}  - {{.Date "previous_period_end" }}) was {{.Number "previous_score"}} avg. calls/day.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) is {{.Number "current_score"}} avg. calls/day.\nThis could be worth looking into.' WHERE type = 'MidtermTrendImmediatelyDown';
					UPDATE message_templates SET template='Your no. of calls per day has been fluctuating since {{.Date "current_period_start" }}.\nThis could be worth looking into.' WHERE type = 'MidtermFluctuationImmediatelyHigh';
					UPDATE message_templates SET template='Fluctuations in the no. of calls per day during the period {{.Date "current_period_start"}} - {{.Date "current_period_end"}} have been successfully stabilized.' WHERE type = 'MidtermFluctuationImmediatelyStabilized';
					UPDATE message_templates SET template='Your avg. no. of calls per day is expected to grow by {{.Number "percentage"}}%% during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}).' WHERE type = 'MidtermPrediction15daysUp';
					UPDATE message_templates SET template='Your avg. no. of calls per day is expected to decline by {{.Number "percentage"}}%% during the period {{.Date "future_period_start" }}  - {{.Date "future_period_start" }}.' WHERE type = 'MidtermPrediction15daysDown';

					UPDATE message_templates SET template='Your avg. objective quality per day has increased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nGreat job!' WHERE type = 'MidtermOQTrend15daysUp';
					UPDATE message_templates SET template='Your avg. objective quality per day has decreased by {{.Number "percentage"}}%%.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThis could be worth looking into.' WHERE type = 'MidtermOQTrend15daysDown';
					UPDATE message_templates SET template='Your avg. objective quality per day has been fluctuating since {{.Date "current_period_start" }}. This could be worth looking into.' WHERE TYPE = 'MidtermOQFluctuationImmediatelyHigh';
					UPDATE message_templates SET template='Fluctuations in avg. objective quality per day during the period {{.Date "current_period_start"}} - {{.Date "current_period_end"}} have been successfully stabilized. Great job!' WHERE type = 'MidtermOQFluctuationImmediatelyStabilized';

					UPDATE message_templates SET template='Your avg. Round Trip Time (RTT) per day has increased.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. RTT/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. RTT/day is {{.Number "current_score"}} ms.\nThis could be worth looking into.' WHERE type = 'MidtermRttTrend15daysUp';
					UPDATE message_templates SET template='Your avg. Round Trip Time (RTT) per day has decreased.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. RTT/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. RTT/day is {{.Number "current_score"}} ms.\nGreat job!' WHERE type = 'MidtermRttTrend15daysDown';
					UPDATE message_templates SET template='Your avg. Round Trip Time (RTT) per day has been fluctuating since {{.Date "current_period_start" }}. This could be worth looking into.' WHERE type = 'MidtermRttFluctuationImmediatelyHigh';
					UPDATE message_templates SET template='Fluctuations in avg. Round Trip Time (RTT) per day during the period {{.Date "current_period_start"}} - {{.Date "current_period_end"}} have been successfully stabilized. Great job!' WHERE type = 'MidtermRttFluctuationImmediatelyStabilized';

					COMMIT;
					`, opts.RootRole))

				return err
			},
		}
	})
}
