package repo

import (
	"context"
	"regexp"
	"testing"
	"time"
	"tradetracker/pkg/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreateTrade(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, db.Close())
	}()
	r, err := NewRepo(WithDB(db))
	require.NoError(t, err)

	trade := &models.Trade{
		InstrumentID: 1,
		Price:        10.0,
		Size:         20,
		Timestamp:    time.Date(2022, time.May, 1, 2, 3, 4, 5, time.UTC),
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		r.queries[createTrade],
	)).WithArgs(trade.InstrumentID, trade.Size, trade.Price, trade.Timestamp).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(1),
	)

	id, err := r.CreateTrade(context.Background(), trade)
	require.NoError(t, err)
	require.Equal(t, 1, id)
}
