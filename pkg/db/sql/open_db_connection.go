package sql

import (
	"github.com/bulutcan99/go-websocket/app/repository"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	AuthQueries *repository.AuthRepo
	UserQueries *repository.UserRepo
}

func SqlQueryInjection() *Queries {
	var db *sqlx.DB

	return &Queries{
		AuthQueries: &repository.AuthRepo{DB: db},
		UserQueries: &repository.UserRepo{DB: db},
	}
}
