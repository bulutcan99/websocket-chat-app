package repository

import (
	"context"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	model2 "github.com/bulutcan99/go-websocket/pkg/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthInterface interface {
	CreateUser(ctx context.Context, u model2.User) error
	GetUserSignByEmail(ctx context.Context, email string) (model2.SignIn, error)
	GetUserById(ctx context.Context, id uuid.UUID) (model2.User, error)
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

func (r *AuthRepo) CreateUser(ctx context.Context, u model2.User) error {
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

func (r *AuthRepo) GetUserSignByEmail(ctx context.Context, email string) (*model2.User, error) {
	var user model2.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.NameSurname, &user.PasswordHash, &user.Status, &user.UserRole)
	if err != nil {
		return &user, custom_error.DatabaseError()
	}

	return &user, nil
}

func (r *AuthRepo) GetUserById(ctx context.Context, id uuid.UUID) (model2.User, error) {
	var user model2.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.NameSurname, &user.PasswordHash, &user.Status, &user.UserRole)
	if err != nil {
		return model2.User{}, custom_error.DatabaseError()
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
