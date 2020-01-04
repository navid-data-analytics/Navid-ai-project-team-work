package migrations

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, opts *migrations.Options) migrations.Migration {
		return migrations.Migration{
			Version: 3,
			Up: func(db migrations.DB) error {
				logger.Info("creating table aid_analytics_states...")
				// Mostly reflects current auth_user table in django based dashboard postgres
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					CREATE TABLE aid_analytics_states(
						id             	SERIAL,
						app_id 			INTEGER NOT NULL,
						keyword   		TEXT NOT NULL,
						data			BYTEA NOT NULL,
						saved_at 		TIMESTAMP WITH TIME ZONE NOT NULL,
						PRIMARY KEY(id),
						CONSTRAINT aid_analytics_states_keyword_idx UNIQUE (app_id, saved_at, keyword)
					);
					GRANT SELECT ON aid_analytics_states TO %s;
					`, opts.RootRole, readRole(opts)))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping table aid_analytics_states...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					DROP TABLE IF EXISTS aid_analytics_states;
				`, opts.RootRole))

				return err
			},
		}
	})
}
