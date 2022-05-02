package apps

import (
	"context"
	"database/sql"

	"tradetracker/internal/pkg/validate"

	"github.com/pkg/errors"
)

// PositionAppCfg configures a PositionApp.
type PositionAppCfg interface {
	ApplyPositionApp(*PositionApp) error
}

// PositionApp is the demo application responsible for carrying out CLI commands.
type PositionApp struct {
	DB *sql.DB `validate:"required"`
}

// NewPositionApp creates a new PositionApp.
func NewPositionApp(cfgs ...PositionAppCfg) (*PositionApp, error) {
	app := &PositionApp{}
	for _, cfg := range cfgs {
		if err := cfg.ApplyPositionApp(app); err != nil {
			return nil, errors.Wrap(err, "apply PositionApp cfg failed")
		}
	}
	if err := validate.Validate().Struct(app); err != nil {
		return nil, errors.Wrap(err, "validate PositionApp failed")
	}
	return app, nil
}

// Run runs the app.
func (app *PositionApp) Run(ctx context.Context, _ []string) error {
	return nil
}
