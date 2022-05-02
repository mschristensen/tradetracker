package repo

import (
	"database/sql"
	"embed"
	"path"
	"tradetracker/internal/pkg/validate"

	"github.com/pkg/errors"
)

//go:embed query/*.sql
var queries embed.FS

// These are query names.
const (
	createTrade = "create_trade.sql"
)

// Repo interacts with the postgres database.
type Repo struct {
	db      *sql.DB           `validate:"required"`
	queries map[string]string `validate:"required"`
}

// NewRepo creates a new Repo for interacting with the database.
// It returns an error if database is not set.
func NewRepo(cfgs ...ConfigFunc) (*Repo, error) {
	r := &Repo{}
	queryFiles := []string{
		createTrade,
		// TODO: add more queries here...
	}
	r.queries = make(map[string]string, len(queryFiles))
	for _, name := range queryFiles {
		q, err := queries.ReadFile(path.Join("query", name))
		if err != nil {
			return nil, errors.Wrap(err, "could not load the query")
		}
		r.queries[name] = string(q)
	}
	for _, cfg := range cfgs {
		cfg(r)
	}
	if err := validate.Validate().Struct(r); err != nil {
		return nil, errors.Wrap(err, "invalid repo")
	}
	return r, nil
}