package position

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
	"tradetracker/pkg/models"

	"github.com/stretchr/testify/require"
)

type tst struct {
	binSize         int64
	initialPosition *models.Position
	trades          []*models.Trade
	positions       []*models.Position
}

var tsts []tst = []tst{
	{
		binSize: 1,
		initialPosition: &models.Position{
			InstrumentID: 1,
			Size:         0,
			Timestamp:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		trades: []*models.Trade{
			{
				InstrumentID: 1,
				Size:         1,
				Timestamp:    time.Date(2022, 1, 1, 0, 0, 1, 0, time.UTC),
			},
			{
				InstrumentID: 1,
				Size:         1,
				Timestamp:    time.Date(2022, 1, 1, 0, 0, 2, 0, time.UTC),
			},
			{
				InstrumentID: 1,
				Size:         1,
				Timestamp:    time.Date(2022, 1, 1, 0, 0, 3, 0, time.UTC),
			},
		},
		positions: []*models.Position{
			{
				InstrumentID: 1,
				Size:         1,
				Timestamp:    time.Date(2022, 1, 1, 0, 0, 1, 0, time.UTC),
			},
			{
				InstrumentID: 1,
				Size:         2,
				Timestamp:    time.Date(2022, 1, 1, 0, 0, 2, 0, time.UTC),
			},
			{
				InstrumentID: 1,
				Size:         3,
				Timestamp:    time.Date(2022, 1, 1, 0, 0, 3, 0, time.UTC),
			},
		},
	},
}

func TestBuilder(t *testing.T) {
	ctx := context.Background()
	for i := range tsts {
		j := i
		t.Run(fmt.Sprintf("test_%d", j), func(t *testing.T) {
			tradesCh := make(chan *models.Trade)
			positionsCh := make(chan *models.Position)
			go func() {
				require.NoError(t,
					NewBinnedBuilder(
						tsts[j].binSize,
						tsts[j].initialPosition,
					).Build(ctx, tradesCh, positionsCh),
				)
			}()
			var actualPositions []*models.Position
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				for pos := range positionsCh {
					actualPositions = append(actualPositions, pos)
				}
			}()
			go func() {
				defer close(tradesCh)
				defer wg.Done()
				for _, trade := range tsts[j].trades {
					tradesCh <- trade
				}
			}()
			wg.Wait()
			for idx := range tsts[j].positions {
				require.Less(t, idx, len(actualPositions))
				require.Equal(t, tsts[j].positions[idx].InstrumentID, actualPositions[idx].InstrumentID, fmt.Sprintf("idx %d", idx))
				require.Equal(t, tsts[j].positions[idx].Size, actualPositions[idx].Size, fmt.Sprintf("idx %d", idx))
				require.Equal(t, tsts[j].positions[idx].Timestamp, actualPositions[idx].Timestamp, fmt.Sprintf("idx %d", idx))
			}
		})
	}
}
