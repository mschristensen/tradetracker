// Package cfg implements functionaltiy to configure an app.
//
// The configuration objects defined here need only be implemented once,
// but can be applied to multiple types.
//
// In order to add support for a new type, the configuration
// need only implement an ApplyX method.
//
package cfg

import (
	"database/sql"

	"tradetracker/internal"
	"tradetracker/internal/app/apps"
	"tradetracker/internal/pkg/db"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var logger logrus.FieldLogger = logrus.StandardLogger()

// DBCfg is configuration for connection to the postgres database using a standard library driver.
type DBCfg struct {
	host, dbName, user, password string
	port                         int
}

// NewDBCfg creates a new DBCfg from the given config.
func NewDBCfg(host string, port int, dbName, user, password string) *DBCfg {
	return &DBCfg{
		host:     host,
		port:     port,
		dbName:   dbName,
		user:     user,
		password: password,
	}
}

// DBFromEnv creates a new DBCfg from the current environment.
func DBFromEnv() *DBCfg {
	return &DBCfg{
		host:     internal.PostgresHost,
		port:     internal.PostgresPort,
		dbName:   internal.PostgresDatabase,
		user:     internal.PostgresUser,
		password: internal.PostgresPassword,
	}
}

func getDBConn(name, host string, port int, database, user, password string) (*sql.DB, error) {
	logger.WithFields(logrus.Fields{
		"host":     host,
		"port":     port,
		"database": database,
		"user":     user,
		"name":     name,
	}).Debugf("Connecting to Postgres %s database...", database)
	dbConn, err := db.NewPGClient(host, port, database, user, password)
	if err != nil {
		return nil, errors.Wrap(err, "new db client failed")
	}
	return dbConn, nil
}

// ApplyTraderApp applies the DBCfg to a CoreApp.
func (cfg DBCfg) ApplyTraderApp(app *apps.TraderApp) error {
	dbConn, err := getDBConn("core", cfg.host, cfg.port, cfg.dbName, cfg.user, cfg.password)
	if err != nil {
		return errors.Wrap(err, "get db conn failed")
	}
	app.DB = dbConn
	return nil
}
