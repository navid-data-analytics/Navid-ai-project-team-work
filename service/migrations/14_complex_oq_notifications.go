package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 14,
			Up: func(db migrations.DB) error {
				logger.Info("adding objective quality complex notification specific messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					INSERT INTO message_templates (type, version, template)
					VALUES
						('MidtermOQCNTrend15daysUpLoss', 1, 'Your avg. objective quality per day has <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe average packet loss which decreased during the same time period contributed to the objective quality improvement.\nGreat job!'),
						('MidtermOQCNTrend15daysDownLoss', 1,'Your avg. objective quality per day has <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe main reason is the average packet loss which increased during the same time period.\nThis could be worth looking into.'),
						('MidtermOQCNTrend15daysUpDelay', 1, 'Your avg. objective quality per day has <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe average round-trip time (RTT) which decreased during the same time period contributed to the objective quality improvement.\nGreat job!'),
						('MidtermOQCNTrend15daysDownDelay', 1, 'Your avg. objective quality per day has <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe main reason is the average round-trip time (RTT) which increased during the same time period.\nThis could be worth looking into.'),
						('MidtermOQCNTrend15daysUpThroughput', 1, 'Your avg. objective quality per day has <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe average throughput which increased during the same time period contributed to the objective quality improvement.\nGreat job!'),
						('MidtermOQCNTrend15daysDownThroughput', 1, 'Your avg. objective quality per day has <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe main reason is the average throughput which decreased during the same time period.\nThis could be worth looking into.'),
						('MidtermOQCNTrend15daysUpLossDelay', 1, 'Your avg. objective quality per day has <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe average packed loss and round-trip time (RTT) which decreased during the same time period contributed to the objective quality improvement.\nGreat job!'),
						('MidtermOQCNTrend15daysDownLossDelay', 1, 'Your avg. objective quality per day has <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe main reasons are the average packet loss and round-trip time (RTT) which increased during the same time period.\nThis could be worth looking into.'),
						('MidtermOQCNTrend15daysUpLossThroughput', 1, 'Your avg. objective quality per day has <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe average packet loss decreased and average throughput increased during the same time period which contributed to the objective quality improvement.\nGreat job!'),
						('MidtermOQCNTrend15daysDownLossThroughput', 1, 'Your avg. objective quality per day has <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe main reasons are the average packet loss increasing and average throughput decreasing during the same time period.\nThis could be worth looking into.'),
						('MidtermOQCNTrend15daysUpThroughputDelay', 1, 'Your avg. objective quality per day has <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe average throughput increased and average round-trip time (RTT) decreased during the same time period which contributed to the objective quality improvement.\nGreat job!'),
						('MidtermOQCNTrend15daysDownThroughputDelay', 1, 'Your avg. objective quality per day has <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe main reasons are the average throughput decreasing and average round-trip time (RTT) increasing during the same time period.\nThis could be worth looking into.'),
						('MidtermOQCNTrend15daysUpAll', 1, 'Your avg. objective quality per day has <span style="color:green; font-weight: bold">increased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe average throughput increased, average round-trip time (RTT) and average packet loss decreased during the same time period which contributed to the objective quality improvement.\nGreat job!'),
						('MidtermOQCNTrend15daysDownAll', 1, 'Your avg. objective quality per day has <span style="color:red; font-weight: bold">decreased</span> by <span style="font-weight: bold">{{.Number "percentage"}}%%</span>.\nPrevious period ({{.Date "previous_period_start" }} - {{.Date "previous_period_end" }}) avg. OQ/day was {{.Number "previous_score"}}.\nCurrent period ({{.Date "current_period_start" }} - {{.Date "current_period_end" }}) avg. OQ/day is {{.Number "current_score"}}.\nThe main reasons are the average throughput decreasing, average round-trip time (RTT) and average packet loss increasing.\nThis could be worth looking into.');
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping objective quality complex notification specific messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					DELETE FROM messages USING message_templates
						WHERE template_id IN (
							SELECT id FROM message_templates
								WHERE type IN (
									'MidtermOQCNTrend15daysUpLoss',
									'MidtermOQCNTrend15daysDownLoss',
									'MidtermOQCNTrend15daysUpDelay',
									'MidtermOQCNTrend15daysDownDelay',
									'MidtermOQCNTrend15daysUpThroughput',
									'MidtermOQCNTrend15daysDownThroughput',
									'MidtermOQCNTrend15daysUpLossDelay',
									'MidtermOQCNTrend15daysDownLossDelay',
									'MidtermOQCNTrend15daysUpLossThroughput',
									'MidtermOQCNTrend15daysDownLossThroughput',
									'MidtermOQCNTrend15daysUpThroughputDelay',
									'MidtermOQCNTrend15daysDownThroughputDelay',
									'MidtermOQCNTrend15daysUpAll',
									'MidtermOQCNTrend15daysDownAll'
								) AND version = 1
						);
					DELETE FROM message_templates
						WHERE type IN (
							'MidtermOQCNTrend15daysUpLoss',
							'MidtermOQCNTrend15daysDownLoss',
							'MidtermOQCNTrend15daysUpDelay',
							'MidtermOQCNTrend15daysDownDelay',
							'MidtermOQCNTrend15daysUpThroughput',
							'MidtermOQCNTrend15daysDownThroughput',
							'MidtermOQCNTrend15daysUpLossDelay',
							'MidtermOQCNTrend15daysDownLossDelay',
							'MidtermOQCNTrend15daysUpLossThroughput',
							'MidtermOQCNTrend15daysDownLossThroughput',
							'MidtermOQCNTrend15daysUpThroughputDelay',
							'MidtermOQCNTrend15daysDownThroughputDelay',
							'MidtermOQCNTrend15daysUpAll',
							'MidtermOQCNTrend15daysDownAll'
						) AND version = 1;
				`, opts.RootRole))

				return err
			},
		}
	})
}
