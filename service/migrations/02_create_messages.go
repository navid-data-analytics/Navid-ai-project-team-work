package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 2,
			Up: func(db migrations.DB) error {
				logger.Info("creating table messages...")
				// Mostly reflects current auth_user table in django based dashboard postgres
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					CREATE TABLE messages(
						id             	SERIAL,
						app_id   		INTEGER NOT NULL,
						template_id		INTEGER NOT NULL,
						generated_at 	TIMESTAMP WITH TIME ZONE NOT NULL,
						data 			BYTEA NOT NULL,
						PRIMARY KEY(id),
						FOREIGN KEY (template_id) REFERENCES message_templates(id)
					);
					GRANT SELECT ON messages TO %s;
					`, opts.RootRole, readRole(opts)))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping table messages...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					DROP TABLE IF EXISTS messages;
				`, opts.RootRole))

				return err
			},
		}
	})
}
