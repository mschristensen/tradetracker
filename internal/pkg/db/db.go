// Package db adds functionality for interacting with the postgres database.
package db

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
)

const (
	maxIdleConn = 80
	maxOpenConn = 80
)

// Schema embeds the schema for the database.
//go:embed schema.sql
var Schema embed.FS

// NewPGClient returns a new postgres connection. You want to share this
// connection throughout the application.
func NewPGClient(host string, port int, dbName, user, password string) (*sql.DB, error) {
	addr := BuildQueryString(user, password, dbName, host, port, "prefer") + " client_encoding=UTF8"
	config, err := pgx.ParseConfig(addr)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	c := stdlib.OpenDB(*config)
	c.SetMaxIdleConns(maxIdleConn)
	c.SetMaxOpenConns(maxOpenConn)
	return c, errors.Wrap(c.Ping(), "pinging database")
}

// BuildQueryString builds a query string.
func BuildQueryString(user, pass, dbname, host string, port int, sslmode string) string {
	parts := []string{}
	if user != "" {
		parts = append(parts, fmt.Sprintf("user=%s", user))
	}
	if pass != "" {
		parts = append(parts, fmt.Sprintf("password=%s", pass))
	}
	if dbname != "" {
		parts = append(parts, fmt.Sprintf("dbname=%s", dbname))
	}
	if host != "" {
		parts = append(parts, fmt.Sprintf("host=%s", host))
	}
	if port != 0 {
		parts = append(parts, fmt.Sprintf("port=%d", port))
	}
	if sslmode != "" {
		parts = append(parts, fmt.Sprintf("sslmode=%s", sslmode))
	}
	return strings.Join(parts, " ")
}
