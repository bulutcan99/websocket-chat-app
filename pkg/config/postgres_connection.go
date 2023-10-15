package config

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/pkg/env"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	DB_MAX_CON           = &env.Env.DbMaxConnections
	DB_MAX_IDLE_CON      = &env.Env.DbMaxIdleConnections
	DB_MAX_LIFE_TIME_CON = &env.Env.DbMaxLifetimeConnections
)

func PostgreSQLConnection() (*sqlx.DB, error) {
	postgresConnURL, err := ConnectionURLBuilder("postgres")
	fmt.Println(postgresConnURL)
	if err != nil {
		return nil, fmt.Errorf("error while trying to connect to the database, %w", err)
	}

	db, err := sqlx.Connect("pgx", postgresConnURL)
	if err != nil {
		return nil, fmt.Errorf("error while trying to connect to the database, %w", err)
	}

	db.SetMaxOpenConns(*DB_MAX_CON)
	db.SetMaxIdleConns(*DB_MAX_IDLE_CON)
	db.SetConnMaxLifetime(time.Duration(*DB_MAX_LIFE_TIME_CON))

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error, not sent ping to the database, %w", err)
	}

	return db, nil
}
