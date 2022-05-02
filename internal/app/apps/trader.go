package apps

import (
	"context"
	"database/sql"
	"io"

	"tradetracker/internal/pkg/pubsub"
	"tradetracker/internal/pkg/repo"
	"tradetracker/internal/pkg/tradeconsumer"
	"tradetracker/internal/pkg/tradesource"
	"tradetracker/internal/pkg/validate"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var logger logrus.FieldLogger = logrus.StandardLogger()

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
func (app *TraderApp) Run(ctx context.Context, _ []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	r, err := repo.NewRepo(repo.WithDB(app.DB))
	if err != nil {
		return errors.Wrap(err, "new repo failed")
	}
	ps := pubsub.NewMemoryPubSub()
	ts := tradesource.NewRandomTradeSource(100, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	tc, err := tradeconsumer.NewTradeConsumer(
		tradeconsumer.WithRepo(r),
		tradeconsumer.WithSubscriber(ps),
	)
	if err != nil {
		return errors.Wrap(err, "new trade consumer failed")
	}
	go func() {
		for {
			trade, err := ts.Next()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				logger.Fatalln(errors.Wrap(err, "next trade failed"))
			}
			if err := ps.Publish(pubsub.Message{
				Topic: pubsub.TradeTopic,
				Value: trade,
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
	return errors.Wrap(tc.Consume(ctx), "consume trades failed")
}
