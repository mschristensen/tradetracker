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

func TestCreatePosition(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		require.NoError(t, db.Close())
	}()
	r, err := NewRepo(WithDB(db))
	require.NoError(t, err)

	position := &models.Position{
		InstrumentID: 1,
		Size:         20,
		Timestamp:    time.Date(2022, time.May, 1, 2, 3, 4, 5, time.UTC),
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		r.queries[createPosition],
	)).WithArgs(position.InstrumentID, position.Size, position.Timestamp.Unix()).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(1),
	)

	id, err := r.CreatePosition(context.Background(), position)
	require.NoError(t, err)
	require.Equal(t, 1, id)
}
