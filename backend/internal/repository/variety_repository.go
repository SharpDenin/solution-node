package repository

import (
	"backend/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type VarietyRepository interface {
	Create(ctx context.Context, variety *models.Variety) error
	GetAll(ctx context.Context) ([]models.Variety, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Variety, error)
	Update(ctx context.Context, variety *models.Variety) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type varietyRepository struct {
	db *DB
}

func NewVarietyRepository(db *DB) VarietyRepository {
	return &varietyRepository{db: db}
}

func (r *varietyRepository) Create(ctx context.Context, v *models.Variety) error {
	query := `
		INSERT INTO varieties (name, description, priority, image_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.Pool.QueryRow(ctx, query,
		v.Name,
		v.Description,
		v.Priority,
		v.ImageURL,
	).Scan(&v.ID, &v.CreatedAt)
}

func (r *varietyRepository) GetAll(ctx context.Context) ([]models.Variety, error) {
	query := `
		SELECT id, name, description, priority, image_url, created_at
		FROM varieties
		ORDER BY name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Variety

	for rows.Next() {
		var v models.Variety

		err := rows.Scan(
			&v.ID,
			&v.Name,
			&v.Description,
			&v.Priority,
			&v.ImageURL,
			&v.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, v)
	}

	return result, rows.Err()
}

func (r *varietyRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Variety, error) {
	query := `
		SELECT id, name, description, priority, image_url, created_at
		FROM varieties
		WHERE id = $1
	`

	var v models.Variety

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&v.ID,
		&v.Name,
		&v.Description,
		&v.Priority,
		&v.ImageURL,
		&v.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("variety not found")
		}
		return nil, err
	}

	return &v, nil
}

func (r *varietyRepository) Update(ctx context.Context, v *models.Variety) error {
	query := `
		UPDATE varieties
		SET name = $1,
		    description = $2,
		    priority = $3,
		    image_url = $4
		WHERE id = $5
	`

	_, err := r.db.Pool.Exec(ctx, query,
		v.Name,
		v.Description,
		v.Priority,
		v.ImageURL,
		v.ID,
	)

	return err
}

func (r *varietyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM varieties
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
