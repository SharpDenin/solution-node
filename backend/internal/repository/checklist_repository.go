package repository

import (
	"backend/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ChecklistRepository interface {
	GetAll(ctx context.Context) ([]models.Checklist, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Checklist, error)
	GetByCode(ctx context.Context, code string) (*models.Checklist, error)
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]models.Checklist, error)
}

type checklistRepository struct {
	db *DB
}

func NewChecklistRepository(db *DB) ChecklistRepository {
	return &checklistRepository{db: db}
}

func (r *checklistRepository) GetAll(ctx context.Context) ([]models.Checklist, error) {
	query := `
		SELECT id, name, code, allowed_role_id, created_at
		FROM checklists
		ORDER BY created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checklists []models.Checklist

	for rows.Next() {
		var checklist models.Checklist

		err := rows.Scan(
			&checklist.ID,
			&checklist.Name,
			&checklist.Code,
			&checklist.AllowedRoleID,
			&checklist.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		checklists = append(checklists, checklist)
	}

	return checklists, rows.Err()
}

func (r *checklistRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Checklist, error) {
	query := `
		SELECT id, name, code, allowed_role_id, created_at
		FROM checklists
		WHERE id = $1
	`

	var checklist models.Checklist

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&checklist.ID,
		&checklist.Name,
		&checklist.Code,
		&checklist.AllowedRoleID,
		&checklist.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("checklist not found")
		}
		return nil, err
	}

	return &checklist, nil
}

func (r *checklistRepository) GetByCode(ctx context.Context, code string) (*models.Checklist, error) {
	query := `
		SELECT id, name, code, allowed_role_id, created_at
		FROM checklists
		WHERE code = $1
	`

	var checklist models.Checklist

	err := r.db.Pool.QueryRow(ctx, query, code).Scan(
		&checklist.ID,
		&checklist.Name,
		&checklist.Code,
		&checklist.AllowedRoleID,
		&checklist.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("checklist not found")
		}
		return nil, err
	}

	return &checklist, nil
}

func (r *checklistRepository) GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]models.Checklist, error) {
	query := `
		SELECT id, name, code, allowed_role_id, created_at
		FROM checklists
		WHERE allowed_role_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checklists []models.Checklist

	for rows.Next() {
		var checklist models.Checklist

		err := rows.Scan(
			&checklist.ID,
			&checklist.Name,
			&checklist.Code,
			&checklist.AllowedRoleID,
			&checklist.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		checklists = append(checklists, checklist)
	}

	return checklists, rows.Err()
}
