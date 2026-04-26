package repository

import (
	"backend/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type QuestionRepository interface {
	Create(ctx context.Context, q *models.Question) error

	GetByID(ctx context.Context, id uuid.UUID) (*models.Question, error)
	GetAll(ctx context.Context) ([]models.Question, error)
	GetByChecklist(ctx context.Context, checklistID uuid.UUID) ([]models.Question, error)

	Update(ctx context.Context, q *models.Question) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type questionRepository struct {
	db *DB
}

func NewQuestionRepository(db *DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(ctx context.Context, q *models.Question) error {
	query := `
		INSERT INTO questions (
			text,
			order_index,
			is_active,
			checklist_id,
			formula,
			image_url
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	return r.db.Pool.QueryRow(ctx, query,
		q.Text,
		q.OrderIndex,
		q.IsActive,
		q.ChecklistID,
		q.Formula,
		q.ImageURL,
	).Scan(&q.ID, &q.CreatedAt)
}

func (r *questionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Question, error) {
	query := `
		SELECT
			id,
			text,
			order_index,
			is_active,
			checklist_id,
			formula,
			image_url,
			created_at
		FROM questions
		WHERE id = $1
	`

	var q models.Question

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&q.ID,
		&q.Text,
		&q.OrderIndex,
		&q.IsActive,
		&q.ChecklistID,
		&q.Formula,
		&q.ImageURL,
		&q.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &q, nil
}

func (r *questionRepository) GetAll(ctx context.Context) ([]models.Question, error) {
	query := `
		SELECT
			id,
			text,
			order_index,
			is_active,
			checklist_id,
			formula,
			image_url,
			created_at
		FROM questions
		ORDER BY order_index ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Question

	for rows.Next() {
		var q models.Question

		err := rows.Scan(
			&q.ID,
			&q.Text,
			&q.OrderIndex,
			&q.IsActive,
			&q.ChecklistID,
			&q.Formula,
			&q.ImageURL,
			&q.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, q)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *questionRepository) GetByChecklist(ctx context.Context, checklistID uuid.UUID) ([]models.Question, error) {
	query := `
		SELECT
			id,
			text,
			order_index,
			is_active,
			checklist_id,
			formula,
			image_url,
			created_at
		FROM questions
		WHERE checklist_id = $1
		  AND is_active = true
		ORDER BY order_index ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, checklistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Question

	for rows.Next() {
		var q models.Question

		err := rows.Scan(
			&q.ID,
			&q.Text,
			&q.OrderIndex,
			&q.IsActive,
			&q.ChecklistID,
			&q.Formula,
			&q.ImageURL,
			&q.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, q)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *questionRepository) Update(ctx context.Context, q *models.Question) error {
	query := `
		UPDATE questions
		SET
			text = $1,
			order_index = $2,
			is_active = $3,
			checklist_id = $4,
			formula = $5,
			image_url = $6
		WHERE id = $7
	`

	_, err := r.db.Pool.Exec(ctx, query,
		q.Text,
		q.OrderIndex,
		q.IsActive,
		q.ChecklistID,
		q.Formula,
		q.ImageURL,
		q.ID,
	)

	return err
}

func (r *questionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE questions
		SET is_active = false
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
