package repository

import (
	"context"
	"github.com/bulutcan99/go-websocket/internal/model"
	config_psql "github.com/bulutcan99/go-websocket/pkg/config/psql"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthInterface interface {
	CreateUser(u model.User) error
	GetUserSignByEmail(email string) (*model.User, error)
	GetUserById(id uuid.UUID) (model.User, error)
	GetUserRoleById(id uuid.UUID) (string, error)
	GetUserEmailById(id uuid.UUID) (string, error)
}

type AuthRepo struct {
	db      *sqlx.DB
	context context.Context
}

func NewAuthUserRepo(psql *config_psql.PostgreSQL) *AuthRepo {
	return &AuthRepo{
		db:      psql.Client,
		context: psql.Context,
	}
}

func (r *AuthRepo) CreateUser(u model.User) error {
	query := `
        INSERT INTO users (
            id,
            uuid,
	          user_name,
	          user_surname,
	          nickname,
            password_hash,
            email,
            user_role,
            user_status,
            created_at,
            updated_at,
            blocked_at       
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(
		r.context,
		query,
		u.Id, u.UUID, u.UserName, u.UserSurName, u.Nickname, u.Passwordhash, u.Email, u.UserRole, u.Status, u.CreatedAt, u.UpdatedAt, u.BlockedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepo) GetUserSignByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.QueryRowContext(r.context, query, email).Scan(&user.Id, &user.UUID, &user.UserName, &user.UserSurName, &user.Nickname, &user.Passwordhash, &user.Email, &user.UserRole, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.BlockedAt)
	if err != nil {
		return &user, custom_error.DatabaseError()
	}

	return &user, nil
}

func (r *AuthRepo) GetUserById(id uuid.UUID) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.QueryRowContext(r.context, query, id).Scan(&user.Id, &user.UUID, &user.UserName, &user.UserSurName, &user.Nickname, &user.Passwordhash, &user.Email, &user.UserRole, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.BlockedAt)
	if err != nil {
		return nil, custom_error.DatabaseError()
	}

	return &user, nil
}

func (r *AuthRepo) GetUserRoleById(id uuid.UUID) (string, error) {
	var userRole string
	query := `SELECT user_role FROM users WHERE id = $1`
	err := r.db.QueryRowContext(r.context, query, id).Scan(&userRole)
	if err != nil {
		return "", custom_error.DatabaseError()
	}

	return userRole, nil
}
