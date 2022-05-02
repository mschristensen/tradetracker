// build +integration
package apps_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"tradetracker/internal"
	"tradetracker/internal/app/apps"
	"tradetracker/internal/app/cfg"
	"tradetracker/pkg/testhelper"

	"github.com/stretchr/testify/require"
)

func TestPositionApp(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip()
	}
	t.Run("TestPosition", testPosition)
}

type positionAppTest struct {
	expectedErr  func(*testing.T, error)
	args         []string
	fixtures     []testhelper.DBFunc
	expectations []testhelper.DBFunc
}

type positionAppTestCfg func(*positionAppTest)

func newPositionAppTest(cfgs ...positionAppTestCfg) *positionAppTest {
	t := &positionAppTest{}
	for _, cfg := range cfgs {
		cfg(t)
	}
	return t
}

func withExpectedError(expectedErr error) positionAppTestCfg {
	return func(t *positionAppTest) {
		t.expectedErr = func(t *testing.T, err error) {
			require.ErrorIs(t, err, expectedErr)
		}
	}
}

func withArgs(args []string) positionAppTestCfg {
	return func(t *positionAppTest) {
		t.args = args
	}
}

func withFixtures(fixtures ...testhelper.DBFunc) positionAppTestCfg {
	return func(t *positionAppTest) {
		t.fixtures = fixtures
	}
}

func withExpectations(expectations ...testhelper.DBFunc) positionAppTestCfg {
	return func(t *positionAppTest) {
		t.expectations = expectations
	}
}

func (tst *positionAppTest) run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	suiteName := strings.Split(t.Name(), "/")[0]
	testName := strings.Split(t.Name(), "/")[1]
	name := strings.ToLower(strings.Join([]string{"tradetracker", suiteName, testName}, "_"))
	dbClient := testhelper.NewDBClient(t,
		name,
		internal.PostgresUser,
		internal.PostgresPassword,
		internal.PostgresHost,
		internal.PostgresPort,
	)
	for _, fixture := range tst.fixtures {
		err := fixture(dbClient)
		require.NoError(t, err)
	}
	app, err := apps.NewPositionApp(
		cfg.NewDBCfg(
			internal.PostgresHost,
			internal.PostgresPort,
			name,
			internal.PostgresUser,
			internal.PostgresPassword,
		),
	)
	require.NoError(t, err)
	err = app.Run(ctx, tst.args)
	if tst.expectedErr == nil {
		require.NoError(t, err)
	} else {
		tst.expectedErr(t, err)
	}
	for _, expectation := range tst.expectations {
		err = expectation(dbClient)
		require.NoError(t, err)
	}
}

func loadSymbols() testhelper.DBFunc {
	return func(db *sql.DB) error {
		_, err := db.Exec(`
			INSERT INTO symbols (symbol) VALUES
				('GBP'),
				('BTC'),
				('ETH'),
				('XRP'),
				('BCH'),
				('LTC');
		`)
		return err
	}
}

func exec(query string) testhelper.DBFunc {
	return func(db *sql.DB) error {
		_, err := db.Exec(query)
		return err
	}
}
