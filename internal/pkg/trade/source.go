package trade

import (
	"context"
	"io"
	"math"
	"math/rand"
	"time"
	"tradetracker/internal/pkg/repo"
	"tradetracker/pkg/models"

	"github.com/pkg/errors"
)

// Source provides an adapter for obtaining trade data.
// It could be anything, from a file, a database, or an external Kafka stream.
type Source interface {
	Prepare(ctx context.Context) error
	Next() (*models.Trade, error)
}

// RandomSource is a source of random trade information.
type RandomSource struct {
	total         int
	num           int
	baseDate      time.Time
	instrumentIDs []int64
	r             *rand.Rand
}

// NewRandomSource creates a new RandomSource.
func NewRandomSource(num int, baseDate time.Time, instrumentIDs []int64) *RandomSource {
	return &RandomSource{
		total:         num,
		baseDate:      baseDate,
		instrumentIDs: instrumentIDs,
	}
}

// Prepare seeds the PRNG.
func (t *RandomSource) Prepare(_ context.Context) error { // nolint:unparam // it's okay that the error is always nil
	t.r = rand.New(rand.NewSource(time.Now().UnixNano())) // nolint:gosec // we don't need a secure PRNG here
	return nil
}

// Next generates a new random trade, or io.EOF if we have generated the maximum number.
func (t *RandomSource) Next() (*models.Trade, error) {
	if t.num >= t.total {
		return nil, io.EOF
	}
	trade := &models.Trade{}
	trade.InstrumentID = t.instrumentIDs[t.r.Intn(len(t.instrumentIDs))]
	trade.Price = math.Round(t.r.Float64()*float64(t.r.Int31n(1000))*100) / 100
	trade.Size = int64(t.r.Int31n(1000))
	// generate random timestamp between baseDate and now
	trade.Timestamp = time.Unix(t.r.Int63n(time.Now().Unix()-t.baseDate.Unix())+t.baseDate.Unix(), 0)
	t.num++
	return trade, nil
}

// RepoSource uses a repo as a source of trade information.
type RepoSource struct {
	repo         repo.TradeRepo
	instrumentID int64
	after        time.Time
	tradeCh      <-chan *models.Trade
}

// NewRepoSource creates a new RepoSource to read trades for the given instrument
// and with a timestamp greater than or equal to `after`.
func NewRepoSource(r repo.TradeRepo, instrumentID int64, after time.Time) *RepoSource {
	return &RepoSource{
		repo:         r,
		instrumentID: instrumentID,
		after:        after,
	}
}

// Prepare initialises the trade channel.
func (t *RepoSource) Prepare(ctx context.Context) error {
	tradeCh, err := t.repo.ReadTrades(ctx, t.instrumentID, t.after)
	if err != nil {
		return errors.Wrap(err, "read trades failed")
	}
	t.tradeCh = tradeCh
	return nil
}

// Next returns the next available trade, or io.EOF if there are no more.
func (t *RepoSource) Next() (*models.Trade, error) {
	if t.tradeCh == nil {
		return nil, io.EOF
	}
	trade, ok := <-t.tradeCh
	if !ok {
		return nil, io.EOF
	}
	return trade, nil
}
