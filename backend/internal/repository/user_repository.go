package repository

import (
	"backend/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User, roleName string) error
	GetByLogin(ctx context.Context, login string) (*models.User, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, string, error)
}

type userRepository struct {
	db *DB
}

func NewUserRepository(db *DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User, roleName string) error {
	query := `
		INSERT INTO users (full_name, login, password_hash, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4))
		RETURNING id, created_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		user.FullName,
		user.Login,
		user.PasswordHash,
		roleName,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errors.New("user already exists")
			}
		}
		return err
	}

	return nil
}

func (r *userRepository) GetByLogin(ctx context.Context, login string) (*models.User, string, error) {
	query := `
		SELECT 
			u.id,
			u.full_name,
			u.login,
			u.password_hash,
			u.role_id,
			u.created_at,
			r.name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.login = $1
	`

	var user models.User
	var roleName string

	err := r.db.Pool.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.FullName,
		&user.Login,
		&user.PasswordHash,
		&user.RoleID,
		&user.CreatedAt,
		&roleName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	return &user, roleName, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, string, error) {
	query := `
		SELECT 
			u.id,
			u.full_name,
			u.login,
			u.password_hash,
			u.role_id,
			u.created_at,
			r.name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`

	var user models.User
	var roleName string

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Login,
		&user.PasswordHash,
		&user.RoleID,
		&user.CreatedAt,
		&roleName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	return &user, roleName, nil
}
