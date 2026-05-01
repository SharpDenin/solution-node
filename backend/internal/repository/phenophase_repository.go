package repository

import (
	"backend/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PhenophaseRepository interface {
	Create(ctx context.Context, phenophase *models.Phenophase) error
	GetAll(ctx context.Context) ([]models.Phenophase, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Phenophase, error)
	Update(ctx context.Context, phenophase *models.Phenophase) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type phenophaseRepository struct {
	db *DB
}

func NewPhenophaseRepository(db *DB) PhenophaseRepository {
	return &phenophaseRepository{db: db}
}

func (r *phenophaseRepository) Create(ctx context.Context, p *models.Phenophase) error {
	query := `
		INSERT INTO phenophases (
			name,
			description,
			image_url,
			order_index,
			min_critical_temperature,
			critical_temperature
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	return r.db.Pool.QueryRow(ctx, query,
		p.Name,
		p.Description,
		p.ImageURL,
		p.OrderIndex,
		p.MinCriticalTemperature,
		p.CriticalTemperature,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *phenophaseRepository) GetAll(ctx context.Context) ([]models.Phenophase, error) {
	query := `
		SELECT 
			id,
			name,
			description,
			image_url,
			order_index,
			min_critical_temperature,
			critical_temperature,
			created_at
		FROM phenophases
		ORDER BY order_index ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Phenophase

	for rows.Next() {
		var p models.Phenophase

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.ImageURL,
			&p.OrderIndex,
			&p.MinCriticalTemperature,
			&p.CriticalTemperature,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, rows.Err()
}

func (r *phenophaseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Phenophase, error) {
	query := `
		SELECT 
			id,
			name,
			description,
			image_url,
			order_index,
			min_critical_temperature,
			critical_temperature,
			created_at
		FROM phenophases
		WHERE id = $1
	`

	var p models.Phenophase

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.ImageURL,
		&p.OrderIndex,
		&p.MinCriticalTemperature,
		&p.CriticalTemperature,
		&p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("phenophase not found")
		}
		return nil, err
	}

	return &p, nil
}

func (r *phenophaseRepository) Update(ctx context.Context, p *models.Phenophase) error {
	query := `
		UPDATE phenophases
		SET name = $1,
		    description = $2,
		    image_url = $3,
		    order_index = $4,
		    min_critical_temperature = $5,
		    critical_temperature = $6
		WHERE id = $7
	`

	commandTag, err := r.db.Pool.Exec(ctx, query,
		p.Name,
		p.Description,
		p.ImageURL,
		p.OrderIndex,
		p.MinCriticalTemperature,
		p.CriticalTemperature,
		p.ID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("phenophase not found")
	}

	return nil
}

func (r *phenophaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM phenophases
		WHERE id = $1
	`

	commandTag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("phenophase not found")
	}

	return nil
}
