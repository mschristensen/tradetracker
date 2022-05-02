package apps

import (
	"context"
	"database/sql"
	"io"
	"strconv"
	"time"

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
func (app *TradeApp) Run(ctx context.Context, args []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// parse the arguments
	if len(args) < 2 {
		return errors.New("requires at least 3 arguments")
	}
	num, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse instrument ID failed")
	}
	instrumentIDs := make([]int64, len(args)-1)
	for i := range args {
		if i == 0 {
			continue
		}
		instrumentID, err := strconv.ParseInt(args[i], 10, 64)
		if err != nil {
			return errors.Wrap(err, "parse instrument ID failed")
		}
		instrumentIDs[i] = instrumentID
	}
	// set up the repository to interact with trades and positions in the database
	r, err := repo.NewRepo(repo.WithDB(app.DB))
	if err != nil {
		return errors.Wrap(err, "new repo failed")
	}
	// create a dummy pubsub stream
	stream := pubsub.NewMemoryPubSub()
	// create a trade source to generate random trade data
	tradeSource := trade.NewRandomSource(
		num,
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), // trades generated from Jan 1, 2000 until now
		instrumentIDs,
	)
	if err := tradeSource.Prepare(ctx); err != nil {
		return errors.Wrap(err, "prepare trade source failed")
	}
	// create a trade processor to process the trade data coming in on the stream
	processor, err := trade.NewProcessor(
		trade.WithRepo(r),
		trade.WithSubscriber(stream),
	)
	if err != nil {
		return errors.Wrap(err, "new trade processor failed")
	}
	// send the random trade data across the stream for it to be processed
	go func() {
		defer func() {
			if err := stream.Close(ctx, pubsub.TradeTopic); err != nil {
				logger.Fatalln(errors.Wrap(err, "close trade stream failed"))
			}
		}()
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
	// process the trade data
	return errors.Wrap(processor.Process(ctx), "process trades failed")
}
