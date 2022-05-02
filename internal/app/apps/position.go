package apps

import (
	"context"
	"database/sql"
	"io"
	"strconv"
	"time"

	"tradetracker/internal/pkg/position"
	"tradetracker/internal/pkg/pubsub"
	"tradetracker/internal/pkg/repo"
	"tradetracker/internal/pkg/trade"
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
func (app *PositionApp) Run(ctx context.Context, args []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// parse the arguments
	if len(args) < 1 {
		return errors.New("missing instrument ID argument")
	}
	instrumentID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse instrument ID failed")
	}
	// set up the repository to interact with trades and positions in the database
	r, err := repo.NewRepo(repo.WithDB(app.DB))
	if err != nil {
		return errors.Wrap(err, "new repo failed")
	}
	// create a dummy pubsub stream
	stream := pubsub.NewMemoryPubSub()
	// create a trade source to read trade data from the repo
	tradeSource := trade.NewRepoSource(r, instrumentID, time.Time{}) // reads all trades, for purposes of this demo
	if err := tradeSource.Prepare(ctx); err != nil {
		return errors.Wrap(err, "prepare trade source failed")
	}
	processor, err := position.NewProcessor(
		position.WithRepo(r),
		position.WithSubscriber(stream),
		position.WithBuilder(
			position.NewBinnedBuilder(1, instrumentID),
		),
	)
	if err != nil {
		return errors.Wrap(err, "new position processor failed")
	}
	// delete positions for the instrument from the date given
	n, err := r.DeletePositions(ctx, instrumentID)
	if err != nil {
		return errors.Wrap(err, "delete positions failed")
	}
	logger.Infof("deleted %d positions", n)
	// send the trade data across the stream for it to be processed
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
	return errors.Wrap(processor.Process(ctx), "process positions failed")
}
