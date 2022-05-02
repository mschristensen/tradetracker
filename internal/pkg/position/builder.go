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
	binWidthSeconds, instrumentID int64
}

// NewBinnedBuilder creates a new BinnedBuilder.
func NewBinnedBuilder(binWidthSeconds, instrumentID int64) *BinnedBuilder {
	return &BinnedBuilder{
		binWidthSeconds: binWidthSeconds,
		instrumentID:    instrumentID,
	}
}

// Build aggregates trades within time windows of binSize seconds to produce positions.
// It assumes that the trades are for a given instrument and are sorted by timestamp; if not, and error is returned.
// NOTE: for now, no actual bin aggregation is implemented, but this is included for demo purposes
func (p *BinnedBuilder) Build(ctx context.Context, in <-chan *models.Trade, out chan<- *models.Position) error {
	defer close(out)
	lastPos := &models.Position{
		InstrumentID: p.instrumentID,
	}
	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "context cancelled")
		case trade, ok := <-in:
			if !ok {
				return nil
			}
			if trade.InstrumentID != lastPos.InstrumentID {
				return ErrInstrumentMismatch
			}
			if trade.Timestamp.Unix() < lastPos.Timestamp.Unix() {
				return errors.Wrapf(
					ErrNotSorted,
					"trade timestamp %s is before position timestamp %s",
					trade.Timestamp.Format(time.RFC3339),
					lastPos.Timestamp.Format(time.RFC3339),
				)
			}
			pos := &models.Position{
				InstrumentID: lastPos.InstrumentID,
				Size:         lastPos.Size + trade.Size,
				Timestamp:    trade.Timestamp,
			}
			out <- pos
			lastPos = pos
		}
	}
}
