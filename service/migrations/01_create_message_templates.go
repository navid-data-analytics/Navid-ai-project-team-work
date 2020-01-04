package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 1,
			Up: func(db migrations.DB) error {
				logger.Info("creating table message_templates...")
				// Mostly reflects current auth_user table in django based dashboard postgres
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					CREATE TABLE message_templates(
						id             	SERIAL,
						type   			TEXT NOT NULL,
						version    		INTEGER NOT NULL,
						template		TEXT NOT NULL,
						created_at 		TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()::timestamp,
						PRIMARY KEY(id)
					);
					CREATE UNIQUE INDEX message_template_versions_idx ON message_templates (type, version);
					GRANT SELECT ON message_templates TO %s;
					`, opts.RootRole, readRole(opts)))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping table message_templates...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					DROP INDEX IF EXISTS message_template_versions_idx;
					DROP TABLE IF EXISTS message_templates;
				`, opts.RootRole))

				return err
			},
		}
	})
}
