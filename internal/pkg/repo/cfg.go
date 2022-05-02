package repo

import (
	"database/sql"
)

// ConfigFunc is used to configure a Repo.
type ConfigFunc func(*Repo)

// WithDB sets the database handle.
func WithDB(db *sql.DB) ConfigFunc {
	return func(r *Repo) {
		r.db = db
	}
}
