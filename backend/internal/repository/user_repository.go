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
	Create(ctx context.Context, user *models.User) error
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type userRepository struct {
	db *DB
}

func NewUserRepository(db *DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (full_name, login, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		user.FullName,
		user.Login,
		user.PasswordHash,
		user.Role,
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

func (r *userRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `
		SELECT id, full_name, login, password_hash, role, created_at
		FROM users
		WHERE login = $1
	`

	var user models.User

	err := r.db.Pool.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.FullName,
		&user.Login,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, full_name, login, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`

	var user models.User

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Login,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
