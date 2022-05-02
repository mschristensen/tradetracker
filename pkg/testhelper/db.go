// Package testhelper implements various test utility functions.
package testhelper

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"tradetracker/internal"
	"tradetracker/internal/pkg/db"

	_ "github.com/jackc/pgx/v4/stdlib" // register pgx database driver
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var migrationsMu sync.Mutex
var mainDBClient *sql.DB

// NewDBClient returns a connection to a test database with the given name. It
// drops the database automatically when the test is finished.
func NewDBClient(t testing.TB, name, owner, pw, host string, port int) *sql.DB {
	t.Helper()
	migrationsMu.Lock()
	defer migrationsMu.Unlock()

	if mainDBClient == nil {
		var err error
		mainDBClient, err = db.NewPGClient(
			internal.PostgresHost,
			internal.PostgresPort,
			internal.PostgresDatabase,
			internal.PostgresUser,
			internal.PostgresPassword,
		)
		if err != nil {
			t.Fatal(errors.Wrap(err, "failed to connect to main database"))
		}
	}

	// use the main database to bootstrap test databases
	name = strings.ToLower(name)
	_, err := mainDBClient.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s' CREATEDB; ALTER USER %s SUPERUSER;", owner, pw, owner))
	if err != nil {
		// only allow an error indicating that the role already exists
		require.Contains(t, err.Error(), fmt.Sprintf(`ERROR: role "%s" already exists (SQLSTATE 42710)`, owner))
	}
	_, err = mainDBClient.Exec(fmt.Sprintf("CREATE DATABASE %s OWNER %s;", name, owner))
	require.NoError(t, err)

	// connect to new database
	dbClient, err := db.NewPGClient(host, port, name, pw, owner)
	require.NoError(t, err)

	schemaSQL, err := db.Schema.ReadFile("schema.sql")
	require.NoError(t, err)
	err = CreateDBSchema(dbClient, string(schemaSQL))
	require.NoError(t, err)

	// teardown databases when done
	t.Cleanup(func() {
		migrationsMu.Lock()
		defer migrationsMu.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		require.NoError(t, dbClient.Close())
		_, err = mainDBClient.ExecContext(ctx, "REVOKE CONNECT ON DATABASE "+name+" FROM public; ")
		require.NoError(t, err)
		_, err = mainDBClient.ExecContext(ctx, `
			SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity
			WHERE pg_stat_activity.datname = $1;
			`, name)
		require.NoError(t, err)
		_, err = mainDBClient.ExecContext(ctx, "DROP DATABASE "+name)
		require.NoError(t, err)
	})
	return dbClient
}

// CreateDBSchema creates the schema and updates the search_path for the given user.
func CreateDBSchema(dbClient *sql.DB, schemaSQL string) error {
	if _, err := dbClient.Exec(schemaSQL); err != nil {
		return errors.Wrap(err, "failed to create schema")
	}
	// Update the search_path as pg_dump resets it by default:
	// https://www.postgresql.org/message-id/ace62b19-f918-3579-3633-b9e19da8b9de%40aklaver.com
	if _, err := dbClient.Exec(`SELECT pg_catalog.set_config('search_path', 'public', false);`); err != nil {
		return errors.Wrap(err, "failed to update search_path")
	}
	return nil
}

// DBFunc supports performing operations on the database.
type DBFunc func(*sql.DB) error
