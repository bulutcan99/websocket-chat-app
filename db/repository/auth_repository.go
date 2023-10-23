package repository

import (
	"context"
	"github.com/bulutcan99/go-websocket/model"
	"github.com/bulutcan99/go-websocket/pkg/config"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthInterface interface {
	CreateUser(u model.User) error
	GetUserSignByEmail(email string) (*model.User, error)
	GetUserById(id uuid.UUID) (model.User, error)
	GetUserRoleById(id uuid.UUID) (string, error)
}

type AuthRepo struct {
	db      *sqlx.DB
	context context.Context
}

func NewAuthUserRepo(psql *config.PostgreSQL) *AuthRepo {
	return &AuthRepo{
		db:      psql.DB,
		context: psql.Context,
	}
}

func (r *AuthRepo) CreateUser(u model.User) error {
	query := `
        INSERT INTO users (
            id,
            created_at,
            updated_at,
            email,
            name_surname,
            password_hash,
            user_status,
            user_role
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(
		r.context,
		query,
		u.ID, u.CreatedAt, u.UpdatedAt, u.Email, u.NameSurname, u.PasswordHash, u.Status, u.UserRole,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepo) GetUserSignByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.QueryRowContext(r.context, query, email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.NameSurname, &user.PasswordHash, &user.Status, &user.UserRole)
	if err != nil {
		return &user, custom_error.DatabaseError()
	}

	return &user, nil
}

func (r *AuthRepo) GetUserById(id uuid.UUID) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.QueryRowContext(r.context, query, id).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.NameSurname, &user.PasswordHash, &user.Status, &user.UserRole)
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
