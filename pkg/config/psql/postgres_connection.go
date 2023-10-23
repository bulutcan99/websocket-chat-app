package psqlconfig

import (
	"context"
	"github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/env"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	doOnce               sync.Once
	client               *sqlx.DB
	DB_MAX_CON           = &env.Env.DbMaxConnections
	DB_MAX_IDLE_CON      = &env.Env.DbMaxIdleConnections
	DB_MAX_LIFE_TIME_CON = &env.Env.DbMaxLifetimeConnections
)

type PostgreSQL struct {
	Client  *sqlx.DB
	Context context.Context
}

func NewPostgreSQLConnection() *PostgreSQL {
	ctx := context.Background()
	postgresConnURL, err := config.ConnectionURLBuilder("postgres")
	if err != nil {
		panic(err)
	}
	doOnce.Do(func() {
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
		client = db
	})

	zap.S().Infof("Connected to Postgres successfully %s", postgresConnURL)
	return &PostgreSQL{
		Client:  client,
		Context: ctx,
	}
}

func (pg *PostgreSQL) Close() {
	err := pg.Client.Close()
	if err != nil {
		zap.S().Errorf("Error while closing the database connection: %s\n", err)
	}

	zap.S().Infof("Connection to Postgres closed successfully")
}
