// build +integration
package apps_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
	"tradetracker/pkg/models"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func testPosition(t *testing.T) {
	t.Parallel()
	now := time.Now()
	newPositionAppTest(
		withArgs([]string{"1"}),
		withFixtures(func(db *sql.DB) error {
			nowCpy := now.Unix()
			for i := 0; i < 5; i++ {
				_, err := db.Exec(`
					INSERT INTO trades (instrument_id, size, price, timestamp)
					VALUES ($1::int, $2::int, $3::numeric, to_timestamp($4::bigint) AT TIME ZONE 'UTC')			
				`, 1, 10, 100.0, nowCpy)
				if err != nil {
					return errors.Wrap(err, "insert trade failed")
				}
				nowCpy += 10 // trades are 10 secs apart
			}
			return nil
		}),
		withExpectations(func(db *sql.DB) error {
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM positions").Scan(&count)
			require.NoError(t, err)
			require.Equal(t, 5, count)
			return nil
		}, func(db *sql.DB) error {
			rows, err := db.Query("SELECT id, instrument_id, size, timestamp FROM positions ORDER BY id")
			require.NoError(t, err)
			i := int64(0)
			for rows.Next() {
				var position models.Position
				require.NoError(t, rows.Scan(
					&position.ID, &position.InstrumentID, &position.Size, &position.Timestamp,
				))
				require.Equal(t, i+1, position.ID, fmt.Sprintf("idx %d", i))
				require.Equal(t, int64(1), position.InstrumentID, fmt.Sprintf("idx %d", i))
				require.Equal(t, 10*(i+1), position.Size, fmt.Sprintf("idx %d", i))
				require.Equal(t, now.Unix()+(i*10), position.Timestamp.Unix(), fmt.Sprintf("idx %d", i))
				i++
			}
			require.NoError(t, rows.Err())
			return nil
		}),
	).run(t)
}
