package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type QuestionFormulaRepository interface {
	GetFormula(ctx context.Context, questionID uuid.UUID, phenophaseID uuid.UUID) (*string, error)
}

type questionFormulaRepository struct {
	db *DB
}

func NewQuestionFormulaRepository(db *DB) QuestionFormulaRepository {
	return &questionFormulaRepository{db: db}
}

func (r *questionFormulaRepository) GetFormula(ctx context.Context, questionID uuid.UUID, phenophaseID uuid.UUID) (*string, error) {
	query := `
		SELECT formula
		FROM question_phenophase_formulas
		WHERE question_id = $1
		  AND phenophase_id = $2
	`

	var formula string

	err := r.db.Pool.QueryRow(ctx, query, questionID, phenophaseID).Scan(&formula)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &formula, nil
}
