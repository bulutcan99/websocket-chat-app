package repository

import (
	"context"
	config_psql "github.com/bulutcan99/go-websocket/pkg/config/psql"
	"github.com/jmoiron/sqlx"
)

type ChatInterface interface {
}

type ChatRepo struct {
	db      *sqlx.DB
	context context.Context
}

func NewChatRepo(psql *config_psql.PostgreSQL) *ChatRepo {
	return &ChatRepo{
		db:      psql.Client,
		context: psql.Context,
	}
}

func (ch *ChatRepo) InsertChatMessage(message []byte) error {
	query := `
				INSERT INTO chat_messages (
						id,
						created_at,
						updated_at,
						message
				) VALUES ($1, $2, $3, $4)`
	_, err := ch.db.ExecContext(
		ch.context,
		query,
		message,
	)
	if err != nil {
		return err
	}

	return nil

}
