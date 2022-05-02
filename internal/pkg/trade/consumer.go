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

// Consumer consumes trades from a pub-sub system and stores them in a repository.
type Consumer struct {
	repo repo.TradeRepo
	sub  pubsub.Subscriber
}

// Cfg is a configuration function for Consumer.
type Cfg func(*Consumer) error

// NewConsumer creates a new Consumer.
func NewConsumer(cfgs ...Cfg) (*Consumer, error) {
	c := &Consumer{}
	for _, cfg := range cfgs {
		if err := cfg(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithRepo sets the repo for the Consumer.
func WithRepo(r repo.TradeRepo) Cfg {
	return func(c *Consumer) error {
		c.repo = r
		return nil
	}
}

// WithSubscriber sets the trade source for the TradeConsumer.
func WithSubscriber(source pubsub.Subscriber) Cfg {
	return func(c *Consumer) error {
		c.sub = source
		return nil
	}
}

// Consume consumes trade messages from the trade source and adds them to the repo.
func (t *Consumer) Consume(ctx context.Context) error {
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
