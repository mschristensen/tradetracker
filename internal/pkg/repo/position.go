package repo

import (
	"context"
	"time"
	"tradetracker/pkg/models"

	"github.com/pkg/errors"
)

// PositionRepo is used to perform CRUD operations on position records in the database.
//go:generate mockery --name PositionRepo --filename position_repo_mock.go
type PositionRepo interface {
	CreatePosition(ctx context.Context, position *models.Position) (int, error)
	ReadPosition(ctx context.Context, instrumentID int64, timestamp time.Time) (*models.Position, error)
}

// CreatePosition creates a new position.
func (r *Repo) CreatePosition(ctx context.Context, position *models.Position) (int, error) {
	var txID int
	if err := r.db.QueryRowContext(ctx,
		r.queries[createPosition],
		position.InstrumentID, position.Size, position.Timestamp.Unix(),
	).Scan(&txID); err != nil {
		return 0, errors.Wrap(err, "could not create position")
	}
	return txID, nil
}

// ReadPosition reads a position for an instrument at a given time.
func (r *Repo) ReadPosition(ctx context.Context, instrumentID int64, timestamp time.Time) (*models.Position, error) {
	var position models.Position
	if err := r.db.QueryRowContext(ctx,
		r.queries[readPosition],
		instrumentID, timestamp.Unix(),
	).Scan(
		&position.ID,
		&position.InstrumentID,
		&position.Size,
		&position.Timestamp,
	); err != nil {
		return nil, errors.Wrap(err, "could not read position")
	}
	return &position, nil
}
