// Package trade implements functionality for interacting with trades.
package trade

import (
	"context"
	"tradetracker/internal/pkg/pubsub"
	"tradetracker/internal/pkg/repo"
	"tradetracker/pkg/models"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var logger logrus.FieldLogger = logrus.StandardLogger()

// Processor consumes trades from a pub-sub system and stores them in a repository.
type Processor struct {
	repo repo.TradeRepo
	sub  pubsub.Subscriber
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
func WithRepo(r repo.TradeRepo) Cfg {
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

// Process consumes trade messages from the trade source and adds them to the repo.
func (t *Processor) Process(ctx context.Context) error {
	err := t.sub.Subscribe(ctx, pubsub.TradeTopic, func(m pubsub.Message) error {
		trade, ok := m.Value.(*models.Trade)
		if !ok {
			return errors.New("could not assert message as trade")
		}
		id, err := t.repo.CreateTrade(ctx, trade)
		if err != nil {
			return errors.Wrap(err, "create trade failed")
		}
		logger.WithFields(logrus.Fields{
			"id":            id,
			"instrument_id": trade.InstrumentID,
			"size":          trade.Size,
			"price":         trade.Price,
			"timestamp":     trade.Timestamp,
		}).Info("added trade")
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "subscribe failed")
	}
	return nil
}
