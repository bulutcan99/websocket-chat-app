package sql

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/jmoiron/sqlx"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // load pgx driver for PostgreSQL
)

func PostgreSQLConnection() (*sqlx.DB, error) {
	Env := env.ParseEnv()
	maxConn := Env.DbMaxConnections
	maxIdleConn := Env.DbMaxIdleConnections
	maxLifeTimeConn := Env.DbMaxLifetimeConnections
	postgresConnURL, err := config.ConnectionURLBuilder("postgres")
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("pgx", postgresConnURL)
	if err != nil {
		return nil, fmt.Errorf("error while trying to connected to database, %w", err)
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(maxLifeTimeConn))

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}
