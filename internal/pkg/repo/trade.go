package repo

import (
	"context"
	"tradetracker/pkg/models"

	"github.com/pkg/errors"
)

// TradeRepo is used to perform CRUD operations on trade records in the database.
//go:generate mockery --name TradeRepo --filename trade_repo_mock.go
type TradeRepo interface {
	CreateTrade(ctx context.Context, trade *models.Trade) (int, error)
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
