package apps

import (
	"context"
	"database/sql"
	"io"

	"tradetracker/internal/pkg/pubsub"
	"tradetracker/internal/pkg/repo"
	"tradetracker/internal/pkg/trade"
	"tradetracker/internal/pkg/validate"

	"github.com/pkg/errors"
)

// TradeAppCfg configures a TradeApp.
type TradeAppCfg interface {
	ApplyTradeApp(*TradeApp) error
}

// TradeApp is the demo application responsible for carrying out CLI commands.
type TradeApp struct {
	DB *sql.DB `validate:"required"`
}

// NewTradeApp creates a new TradeApp.
func NewTradeApp(cfgs ...TradeAppCfg) (*TradeApp, error) {
	app := &TradeApp{}
	for _, cfg := range cfgs {
		if err := cfg.ApplyTradeApp(app); err != nil {
			return nil, errors.Wrap(err, "apply TradeApp cfg failed")
		}
	}
	if err := validate.Validate().Struct(app); err != nil {
		return nil, errors.Wrap(err, "validate TradeApp failed")
	}
	return app, nil
}

// Run runs the app.
func (app *TradeApp) Run(ctx context.Context, _ []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	r, err := repo.NewRepo(repo.WithDB(app.DB))
	if err != nil {
		return errors.Wrap(err, "new repo failed")
	}
	ps := pubsub.NewMemoryPubSub()
	ts := trade.NewRandomSource(100, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	tc, err := trade.NewConsumer(
		trade.WithRepo(r),
		trade.WithSubscriber(ps),
	)
	if err != nil {
		return errors.Wrap(err, "new trade consumer failed")
	}
	go func() {
		for {
			tr, err := ts.Next()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				logger.Fatalln(errors.Wrap(err, "next trade failed"))
			}
			if err := ps.Publish(pubsub.Message{
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
	return errors.Wrap(tc.Consume(ctx), "consume trades failed")
}
