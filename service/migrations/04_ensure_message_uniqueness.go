package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 4,
			Up: func(db migrations.DB) error {
				logger.Info("adding uniqueness constraint to messages...")
				// Mostly reflects current auth_user table in django based dashboard postgres
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					ALTER TABLE messages ADD CONSTRAINT message_uniqueness_idx UNIQUE(app_id, generated_at, template_id);
					`, opts.RootRole))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping adding uniqueness constraint to messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					ALTER TABLE messages DROP CONSTRAINT message_uniqueness_idx;
				`, opts.RootRole))

				return err
			},
		}
	})
}
