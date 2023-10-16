package repository

import (
	"context"
	"github.com/bulutcan99/go-websocket/app/model"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthInterface interface {
	CreateUser(ctx context.Context, u model.User) error
	GetUserById(ctx context.Context, id uuid.UUID) (model.User, error)
	GetUserRoleById(ctx context.Context, id uuid.UUID) (string, error)
}

type AuthRepo struct {
	DB *sqlx.DB
}

func NewAuthUserRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{
		DB: db,
	}
}

func (r *AuthRepo) CreateUser(ctx context.Context, u model.User) error {
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
	_, err := r.DB.ExecContext(
		ctx,
		query,
		u.ID, u.CreatedAt, u.UpdatedAt, u.Email, u.NameSurname, u.PasswordHash, u.Status, u.UserRole,
	)
	if err != nil {
		return custom_error.DatabaseError()
	}

	return nil
}

func (r *AuthRepo) GetUserById(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.NameSurname, &user.PasswordHash, &user.Status, &user.UserRole)
	if err != nil {
		return model.User{}, custom_error.DatabaseError()
	}

	return user, nil
}

func (r *AuthRepo) GetUserRoleById(ctx context.Context, id uuid.UUID) (string, error) {
	var userRole string
	query := `SELECT user_role FROM users WHERE id = $1`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&userRole)
	if err != nil {
		return "", custom_error.DatabaseError()
	}

	return userRole, nil
}
