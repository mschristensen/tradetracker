package apps

import (
	"context"
	"database/sql"
	"io"
	"time"

	"tradetracker/internal/pkg/position"
	"tradetracker/internal/pkg/pubsub"
	"tradetracker/internal/pkg/repo"
	"tradetracker/internal/pkg/trade"
	"tradetracker/internal/pkg/validate"
	"tradetracker/pkg/models"

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
	r, err := repo.NewRepo(repo.WithDB(app.DB))
	if err != nil {
		return errors.Wrap(err, "new repo failed")
	}
	stream := pubsub.NewMemoryPubSub()
	tradeSource := trade.NewRandomSource(100, []int64{1}) // TODO use a sorted trade source for a given instrument
	processor, err := position.NewProcessor(
		position.WithRepo(r),
		position.WithSubscriber(stream),
		position.WithPositionBuilder(
			position.NewBinnedPositionBuilder(1000, &models.Position{
				InstrumentID: 1,
				Size:         0,
				Timestamp:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			}),
		),
	)
	if err != nil {
		return errors.Wrap(err, "new position processor failed")
	}
	go func() {
		for {
			tr, err := tradeSource.Next()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				logger.Fatalln(errors.Wrap(err, "next trade failed"))
			}
			if err := stream.Publish(pubsub.Message{
				Topic: pubsub.TradeTopic,
				Value: tr,
			}); err != nil {
				logger.Fatalln(errors.Wrap(err, "publish trade failed"))
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
	return errors.Wrap(processor.Process(ctx), "process trades failed")
}
