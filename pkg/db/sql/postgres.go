package sql

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/jmoiron/sqlx"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	DB_MAX_CON           = &env.Env.DbMaxConnections
	DB_MAX_IDLE_CON      = &env.Env.DbMaxIdleConnections
	DB_MAX_LIFE_TIME_CON = &env.Env.DbMaxLifetimeConnections
)

func PostgreSQLConnection() (*sqlx.DB, error) {
	maxConn := *DB_MAX_CON
	maxIdleConn := *DB_MAX_IDLE_CON
	maxLifeTimeConn := *DB_MAX_LIFE_TIME_CON
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
