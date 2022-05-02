// Package position implements functionality for interacting with positions.
package position

import (
	"context"
	"sync"
	"tradetracker/internal/pkg/pubsub"
	"tradetracker/internal/pkg/repo"
	"tradetracker/pkg/models"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var logger logrus.FieldLogger = logrus.StandardLogger()

// Processor aggregates trades from a pub-sub system to build positions and stores them in a repository.
type Processor struct {
	repo    repo.PositionRepo
	sub     pubsub.Subscriber
	builder Builder
}

// Cfg is a configuration function for Processor.
type Cfg func(*Processor) error

// NewProcessor creates a new Processor.
func NewProcessor(cfgs ...Cfg) (*Processor, error) {
	c := &Processor{}
	for _, cfg := range cfgs {
		if err := cfg(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithRepo sets the repo for the Processor.
func WithRepo(r repo.PositionRepo) Cfg {
	return func(c *Processor) error {
		c.repo = r
		return nil
	}
}

// WithSubscriber sets the trade source for the Processor.
func WithSubscriber(source pubsub.Subscriber) Cfg {
	return func(c *Processor) error {
		c.sub = source
		return nil
	}
}

// WithBuilder sets the position builder for the Processor.
func WithBuilder(builder Builder) Cfg {
	return func(c *Processor) error {
		c.builder = builder
		return nil
	}
}

// Process consumes trade messages from the trade source and uses them to build positions.
func (t *Processor) Process(ctx context.Context) error {
	tradeCh := make(chan *models.Trade)
	positionCh := make(chan *models.Position)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := t.builder.Build(ctx, tradeCh, positionCh); err != nil {
			logger.Fatalln(errors.Wrap(err, "build positions failed"))
		}
	}()
	go func() {
		defer wg.Done()
		for pos := range positionCh {
			id, err := t.repo.CreatePosition(ctx, pos)
			if err != nil {
				logger.Fatalln(errors.Wrap(err, "create position failed"))
			}
			logger.WithFields(logrus.Fields{
				"id":            id,
				"instrument_id": pos.InstrumentID,
				"size":          pos.Size,
				"timestamp":     pos.Timestamp,
			}).Info("added position")
		}
	}()
	err := t.sub.Subscribe(ctx, pubsub.TradeTopic, func(m pubsub.Message) error {
		trade, ok := m.Value.(*models.Trade)
		if !ok {
			return errors.New("could not assert message as trade")
		}
		tradeCh <- trade
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "subscribe failed")
	}
	close(tradeCh)
	wg.Wait()
	return nil
}
