package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 11,
			Up: func(db migrations.DB) error {
				logger.Info("remove all message data and state...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					BEGIN TRANSACTION;
					DELETE FROM messages;
					DELETE FROM aid_analytics_states;
					COMMIT;
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("cant revert 'remove all message data and state'...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					`, opts.RootRole))

				return err
			},
		}
	})
}
