package apps

import (
	"context"
	"database/sql"
	"time"

	"tradetracker/internal/pkg/repo"
	"tradetracker/internal/pkg/validate"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// QueryAppCfg configures a QueryApp.
type QueryAppCfg interface {
	ApplyQueryApp(*QueryApp) error
}

// QueryApp is the demo application responsible for carrying out CLI commands.
type QueryApp struct {
	DB *sql.DB `validate:"required"`
}

// NewQueryApp creates a new QueryApp.
func NewQueryApp(cfgs ...QueryAppCfg) (*QueryApp, error) {
	app := &QueryApp{}
	for _, cfg := range cfgs {
		if err := cfg.ApplyQueryApp(app); err != nil {
			return nil, errors.Wrap(err, "apply QueryApp cfg failed")
		}
	}
	if err := validate.Validate().Struct(app); err != nil {
		return nil, errors.Wrap(err, "validate QueryApp failed")
	}
	return app, nil
}

// Run runs the app.
func (app *QueryApp) Run(ctx context.Context, _ []string) error {
	r, err := repo.NewRepo(repo.WithDB(app.DB))
	if err != nil {
		return errors.Wrap(err, "new repo failed")
	}
	instrumentID := int64(1)
	timestamp, err := time.Parse(time.RFC3339, "2022-05-02T18:33:36Z")
	if err != nil {
		return errors.Wrap(err, "parse timestamp failed")
	}
	pos, err := r.ReadPosition(ctx, instrumentID, timestamp)
	if err != nil {
		return errors.Wrap(err, "read position failed")
	}
	logger.WithFields(logrus.Fields{
		"instrument_id": pos.InstrumentID,
		"size":          pos.Size,
		"timestamp":     pos.Timestamp,
	}).Info("position found")
	return nil
}
