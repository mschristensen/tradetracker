package position

import (
	"context"
	"time"
	"tradetracker/pkg/models"

	"github.com/pkg/errors"
)

// Builder is used to build positions from trades.
type Builder interface {
	Build(ctx context.Context, in <-chan *models.Trade, out chan<- *models.Position) error
}

// BinnedBuilder builds positions from trades that occur within the a fixed-width bin.
type BinnedBuilder struct {
	binWidthSeconds int64
	initialPosition *models.Position
}

// NewBinnedBuilder creates a new BinnedBuilder.
func NewBinnedBuilder(binWidthSeconds int64, initialPosition *models.Position) *BinnedBuilder {
	return &BinnedBuilder{
		binWidthSeconds: binWidthSeconds,
		initialPosition: initialPosition,
	}
}

// Build aggregates trades within time windows of binSize seconds to produce positions.
// It assumes that the trades are for a given instrument and are sorted by timestamp; if not, and error is returned.
func (p *BinnedBuilder) Build(ctx context.Context, in <-chan *models.Trade, out chan<- *models.Position) error {
	defer close(out)
	pos := p.initialPosition
	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "context cancelled")
		case trade, ok := <-in:
			if !ok {
				return nil
			}
			if trade.InstrumentID != pos.InstrumentID {
				return ErrInstrumentMismatch
			}
			if trade.Timestamp.Unix() < pos.Timestamp.Unix() {
				return errors.Wrapf(
					ErrNotSorted,
					"trade timestamp %s is before position timestamp %s",
					trade.Timestamp.Format(time.RFC3339),
					pos.Timestamp.Format(time.RFC3339),
				)
			}
			if trade.Timestamp.Unix()-pos.Timestamp.Unix() <= p.binWidthSeconds {
				pos.Size += trade.Size
				continue
			}
			pos.Timestamp = trade.Timestamp
			out <- pos
		}
	}
}
