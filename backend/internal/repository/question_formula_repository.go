package repository

import (
	"backend/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type QuestionFormulaRepository interface {
	CreateOrUpdate(ctx context.Context, formula *models.QuestionPhenophaseFormula) error

	GetByQuestion(ctx context.Context, questionID uuid.UUID) ([]models.QuestionPhenophaseFormula, error)
	GetFormula(ctx context.Context, questionID uuid.UUID, phenophaseID uuid.UUID) (*string, error)

	DeleteByQuestion(ctx context.Context, questionID uuid.UUID) error
	DeleteByQuestionAndPhenophase(ctx context.Context, questionID uuid.UUID, phenophaseID uuid.UUID) error
}

type questionFormulaRepository struct {
	db *DB
}

func NewQuestionFormulaRepository(db *DB) QuestionFormulaRepository {
	return &questionFormulaRepository{db: db}
}

func (r *questionFormulaRepository) CreateOrUpdate(ctx context.Context, formula *models.QuestionPhenophaseFormula) error {
	query := `
		INSERT INTO question_phenophase_formulas (
			question_id,
			phenophase_id,
			formula
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (question_id, phenophase_id)
		DO UPDATE SET
			formula = EXCLUDED.formula,
			updated_at = now()
		RETURNING id, created_at, updated_at
	`

	return r.db.Pool.QueryRow(ctx, query,
		formula.QuestionID,
		formula.PhenophaseID,
		formula.Formula,
	).Scan(
		&formula.ID,
		&formula.CreatedAt,
		&formula.UpdatedAt,
	)
}

func (r *questionFormulaRepository) GetByQuestion(ctx context.Context, questionID uuid.UUID) ([]models.QuestionPhenophaseFormula, error) {
	query := `
		SELECT
			id,
			question_id,
			phenophase_id,
			formula,
			created_at,
			updated_at
		FROM question_phenophase_formulas
		WHERE question_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.QuestionPhenophaseFormula

	for rows.Next() {
		var formula models.QuestionPhenophaseFormula

		err := rows.Scan(
			&formula.ID,
			&formula.QuestionID,
			&formula.PhenophaseID,
			&formula.Formula,
			&formula.CreatedAt,
			&formula.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, formula)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
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

func (r *questionFormulaRepository) DeleteByQuestion(ctx context.Context, questionID uuid.UUID) error {
	query := `
		DELETE FROM question_phenophase_formulas
		WHERE question_id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, questionID)
	return err
}

func (r *questionFormulaRepository) DeleteByQuestionAndPhenophase(
	ctx context.Context,
	questionID uuid.UUID,
	phenophaseID uuid.UUID,
) error {
	query := `
		DELETE FROM question_phenophase_formulas
		WHERE question_id = $1
		  AND phenophase_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, questionID, phenophaseID)
	return err
}
