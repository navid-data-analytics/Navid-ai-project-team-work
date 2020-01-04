package migrations

import "github.com/callstats-io/go-common/postgres/migrations"

// MetaKeyReadRole is the key under which db access role should be found
const MetaKeyReadRole = "readRole"

func readRole(opts *migrations.Options) string {
	if opts.Meta[MetaKeyReadRole] != "" {
		return opts.Meta[MetaKeyReadRole].(string)
	}
	return opts.RootRole
}
