package config

import (
	"context"
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

type PostgreSQL struct {
	DB      *sqlx.DB
	Context context.Context
}

func NewPostgreSQLConnection() *PostgreSQL {
	postgresConnURL, err := ConnectionURLBuilder("postgres")
	fmt.Println(postgresConnURL)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	db, err := sqlx.ConnectContext(ctx, "pgx", postgresConnURL)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(*DB_MAX_CON)
	db.SetMaxIdleConns(*DB_MAX_IDLE_CON)
	db.SetConnMaxLifetime(time.Duration(*DB_MAX_LIFE_TIME_CON))
	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	return &PostgreSQL{DB: db, Context: ctx}
}

func (pg *PostgreSQL) Close() {
	err := pg.DB.Close()
	if err != nil {
		fmt.Printf("Error while closing the database connection: %s\n", err)
	}
}
