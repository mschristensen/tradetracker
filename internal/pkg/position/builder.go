package position

import (
	"context"
	"tradetracker/pkg/models"

	"github.com/pkg/errors"
)

// PositionBuilder is used to build positions from trades.
type PositionBuilder interface {
	Build(ctx context.Context, in <-chan *models.Trade, out chan<- *models.Position) error
}

// BinnedPositionBuilder builds positions from trades that occur within the a fixed-width bin.
type BinnedPositionBuilder struct {
	binWidthSeconds int64
	initialPosition *models.Position
}

// NewBinnedPositionBuilder creates a new BinnedPositionBuilder.
func NewBinnedPositionBuilder(binWidthSeconds int64, initialPosition *models.Position) *BinnedPositionBuilder {
	return &BinnedPositionBuilder{
		binWidthSeconds: binWidthSeconds,
		initialPosition: initialPosition,
	}
}

// Build aggregates trades within time windows of binSize seconds to produce positions.
// It assumes that the trades are for a given instrument and are sorted by timestamp; if not, and error is returned.
func (p *BinnedPositionBuilder) Build(ctx context.Context, in <-chan *models.Trade, out chan<- *models.Position) error {
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
				return ErrNotSorted
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
