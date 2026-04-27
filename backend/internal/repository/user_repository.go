package repository

import (
	"backend/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User, roleName string) error

	GetByLogin(ctx context.Context, login string) (*models.User, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, string, error)

	GetAll(ctx context.Context) ([]models.User, map[uuid.UUID]string, error)

	Update(ctx context.Context, user *models.User, roleName string) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	Block(ctx context.Context, id uuid.UUID) error
	Unblock(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *DB
}

func NewUserRepository(db *DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User, roleName string) error {
	query := `
		INSERT INTO users (full_name, login, password_hash, role_id, position)
		SELECT $1, $2, $3, roles.id, $5
		FROM roles
		WHERE roles.name = $4
		RETURNING id, role_id, is_blocked, is_deleted, created_at, updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		user.FullName,
		user.Login,
		user.PasswordHash,
		roleName,
		user.Position,
	).Scan(
		&user.ID,
		&user.RoleID,
		&user.IsBlocked,
		&user.IsDeleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errors.New("user already exists")
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("invalid role")
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
			u.position,
			u.is_blocked,
			u.is_deleted,
			u.created_at,
			u.updated_at,
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
		&user.Position,
		&user.IsBlocked,
		&user.IsDeleted,
		&user.CreatedAt,
		&user.UpdatedAt,
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
			u.position,
			u.is_blocked,
			u.is_deleted,
			u.created_at,
			u.updated_at,
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
		&user.Position,
		&user.IsBlocked,
		&user.IsDeleted,
		&user.CreatedAt,
		&user.UpdatedAt,
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

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, map[uuid.UUID]string, error) {
	query := `
		SELECT 
			u.id,
			u.full_name,
			u.login,
			u.password_hash,
			u.role_id,
			u.position,
			u.is_blocked,
			u.is_deleted,
			u.created_at,
			u.updated_at,
			r.name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE NOT (u.login = 'admin' AND r.name = 'admin')
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	rolesByUserID := make(map[uuid.UUID]string)

	for rows.Next() {
		var user models.User
		var roleName string

		err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Login,
			&user.PasswordHash,
			&user.RoleID,
			&user.Position,
			&user.IsBlocked,
			&user.IsDeleted,
			&user.CreatedAt,
			&user.UpdatedAt,
			&roleName,
		)
		if err != nil {
			return nil, nil, err
		}

		users = append(users, user)
		rolesByUserID[user.ID] = roleName
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return users, rolesByUserID, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User, roleName string) error {
	query := `
		UPDATE users
		SET 
			full_name = $1,
			login = $2,
			role_id = roles.id,
			position = $4,
			updated_at = now()
		FROM roles
		WHERE users.id = $5
		  AND roles.name = $3
		RETURNING users.role_id, users.updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		user.FullName,
		user.Login,
		roleName,
		user.Position,
		user.ID,
	).Scan(&user.RoleID, &user.UpdatedAt)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errors.New("user login already exists")
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user or role not found")
		}

		return err
	}

	return nil
}

func (r *userRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET 
			is_deleted = true,
			is_blocked = true,
			updated_at = now()
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *userRepository) Restore(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET 
			is_deleted = false,
			is_blocked = false,
			updated_at = now()
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *userRepository) Block(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET 
			is_blocked = true,
			updated_at = now()
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *userRepository) Unblock(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET 
			is_blocked = false,
			updated_at = now()
		WHERE id = $1
		  AND is_deleted = false
	`

	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
