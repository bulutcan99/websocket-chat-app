package repository

import (
	"context"
	"github.com/bulutcan99/go-websocket/model"
	config_psql "github.com/bulutcan99/go-websocket/pkg/config/psql"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/bulutcan99/go-websocket/pkg/utility"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserInterface interface {
	GetUserSelf(id uuid.UUID) (*model.User, error)
	GetShownedUserByEmail(email string) (*model.UserShown, error)
	UpdatePassword(id uuid.UUID, oldPassword string, newPassword string) error
}

type UserRepo struct {
	db      *sqlx.DB
	context context.Context
}

func NewUserRepo(psql *config_psql.PostgreSQL) *UserRepo {
	return &UserRepo{
		db:      psql.Client,
		context: psql.Context,
	}
}

func (r *UserRepo) GetUserSelf(id uuid.UUID) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.QueryRowContext(r.context, query, id).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.NameSurname, &user.PasswordHash, &user.Status, &user.UserRole)
	if err != nil {
		return nil, custom_error.DatabaseError()
	}

	return &user, nil
}

func (r *UserRepo) GetShownedUserByEmail(email string) (*model.UserShown, error) {
	var user model.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.QueryRowContext(r.context, query, email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.NameSurname, &user.PasswordHash, &user.Status, &user.UserRole)
	if err != nil {
		return nil, custom_error.DatabaseError()
	}

	return &model.UserShown{
		Email:       user.Email,
		NameSurname: user.NameSurname,
		CreatedAt:   user.CreatedAt,
		UserRole:    user.UserRole,
		Status:      user.Status,
	}, nil
}

func (r *UserRepo) UpdatePassword(id uuid.UUID, newPassword string) error {
	var user model.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.QueryRowContext(r.context, query, id).Scan(&user.PasswordHash)
	if err != nil {
		return custom_error.DatabaseError()
	}

	updateQuery := `UPDATE users SET password_hash = $1 WHERE id = $2`
	hashedPassword := utility.GeneratePassword(newPassword)
	_, updateError := r.db.ExecContext(
		r.context,
		updateQuery,
		hashedPassword)

	if updateError != nil {
		return err
	}

	return nil
}
