package apps

import (
	"context"
	"database/sql"

	"tradetracker/internal/pkg/validate"

	"github.com/pkg/errors"
)

// TraderAppCfg configures a TraderApp.
type TraderAppCfg interface {
	ApplyTraderApp(*TraderApp) error
}

// TraderApp is the demo application responsible for carrying out CLI commands.
type TraderApp struct {
	DB *sql.DB `validate:"required"`
}

// NewTraderApp creates a new TraderApp.
func NewTraderApp(cfgs ...TraderAppCfg) (*TraderApp, error) {
	app := &TraderApp{}
	for _, cfg := range cfgs {
		if err := cfg.ApplyTraderApp(app); err != nil {
			return nil, errors.Wrap(err, "apply TraderApp cfg failed")
		}
	}
	if err := validate.Validate().Struct(app); err != nil {
		return nil, errors.Wrap(err, "validate TraderApp failed")
	}
	return app, nil
}

// Run runs the app.
func (app *TraderApp) Run(ctx context.Context, args []string) error {
	return nil
}
