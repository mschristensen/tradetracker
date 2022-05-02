package repo

import (
	"context"
	"time"
	"tradetracker/pkg/models"

	"github.com/pkg/errors"
)

// TradeRepo is used to perform CRUD operations on trade records in the database.
//go:generate mockery --name TradeRepo --filename trade_repo_mock.go
type TradeRepo interface {
	CreateTrade(ctx context.Context, trade *models.Trade) (int, error)
	ReadTrades(ctx context.Context, instrumentID int64, after time.Time) (<-chan *models.Trade, error)
}

// CreateTrade creates a new trade.
func (r *Repo) CreateTrade(ctx context.Context, trade *models.Trade) (int, error) {
	var txID int
	if err := r.db.QueryRowContext(ctx,
		r.queries[createTrade],
		trade.InstrumentID, trade.Size, trade.Price, trade.Timestamp.Unix(),
	).Scan(&txID); err != nil {
		return 0, errors.Wrap(err, "could not create trade")
	}
	return txID, nil
}

// ReadTrades reads trades from the database and sends them on the returned channel.
func (r *Repo) ReadTrades(ctx context.Context, instrumentID int64, after time.Time) (<-chan *models.Trade, error) { // nolint:unparam // it's okay that the error is always nil
	ch := make(chan *models.Trade)
	go func() {
		rows, err := r.db.QueryContext(ctx, r.queries[readTrades], instrumentID, after.Unix())
		if err != nil {
			logger.Fatalln(errors.Wrap(err, "could not read trades"))
		}
		defer close(ch)
		defer func() {
			if err := rows.Close(); err != nil {
				logger.Fatalln(errors.Wrap(err, "close rows failed"))
			}
		}()
		for rows.Next() {
			var trade models.Trade
			if err := rows.Scan(
				&trade.ID,
				&trade.InstrumentID,
				&trade.Price,
				&trade.Size,
				&trade.Timestamp,
			); err != nil {
				logger.Fatalln(errors.Wrap(err, "scan failed"))
				continue
			}
			ch <- &trade
		}
		if err := rows.Err(); err != nil {
			logger.Fatalln(errors.Wrap(err, "rows failed"))
		}
	}()
	return ch, nil
}
