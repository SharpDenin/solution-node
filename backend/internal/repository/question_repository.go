package repository

import (
	"backend/internal/models"
	"context"
)

type QuestionRepository interface {
	Create(ctx context.Context, q *models.Question) error
	GetAll(ctx context.Context) ([]models.Question, error)
	Update(ctx context.Context, q *models.Question) error
	Delete(ctx context.Context, id string) error
}

type questionRepository struct {
	db *DB
}

func NewQuestionRepository(db *DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(ctx context.Context, q *models.Question) error {

	query := `
		INSERT INTO questions (text, order_index, is_active)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return r.db.Pool.QueryRow(ctx, query,
		q.Text,
		q.OrderIndex,
		q.IsActive,
	).Scan(&q.ID, &q.CreatedAt)
}

func (r *questionRepository) GetAll(ctx context.Context) ([]models.Question, error) {

	query := `
		SELECT id, text, order_index, is_active, created_at
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
			&q.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, q)
	}

	return result, nil
}

func (r *questionRepository) Update(ctx context.Context, q *models.Question) error {

	query := `
		UPDATE questions
		SET text = $1,
		    order_index = $2,
		    is_active = $3
		WHERE id = $4
	`

	_, err := r.db.Pool.Exec(ctx, query,
		q.Text,
		q.OrderIndex,
		q.IsActive,
		q.ID,
	)

	return err
}

func (r *questionRepository) Delete(ctx context.Context, id string) error {

	query := `
		UPDATE questions
		SET is_active = false
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
