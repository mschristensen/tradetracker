// Package tradesource provides an adapter for obtaining trade data.
package tradesource

import (
	"io"
	"math"
	"math/rand"
	"time"
	"tradetracker/pkg/models"
)

// TradeSource is a source of trade information.
// It could be anything, from a file, a database, or an external Kafka stream.
type TradeSource interface {
	Next() (*models.Trade, error)
}

// RandomTradeSource is a source of random trade information.
type RandomTradeSource struct {
	total         int
	n             int
	r             *rand.Rand
	instrumentIDs []int64
}

// NewRandomTradeSource creates a new RandomTradeSource.
func NewRandomTradeSource(num int, instrumentIDs []int64) *RandomTradeSource {
	return &RandomTradeSource{
		total:         num,
		instrumentIDs: instrumentIDs,
		r:             rand.New(rand.NewSource(time.Now().UnixNano())), // nolint:gosec // we don't need a secure PRNG here
	}
}

// Next returns the next trade, or io.EOF if we have generated the maximum number.
func (t *RandomTradeSource) Next() (*models.Trade, error) {
	if t.n >= t.total {
		return nil, io.EOF
	}
	trade := &models.Trade{}
	trade.InstrumentID = t.instrumentIDs[t.r.Intn(len(t.instrumentIDs))]
	trade.Price = math.Round(t.r.Float64()*float64(t.r.Int31n(1000))*100) / 100
	trade.Size = int64(t.r.Int31n(1000))
	baseDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()       // earliest date we will generate
	trade.Timestamp = time.Unix(t.r.Int63n(time.Now().Unix()-baseDate)+baseDate, 0) // generate random timestamp between baseDate and now
	t.n++
	return trade, nil
}
