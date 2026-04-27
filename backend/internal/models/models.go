package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	FullName     string    `db:"full_name"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	RoleID       uuid.UUID `db:"role_id"`
	Position     *string   `db:"position"`
	IsBlocked    bool      `db:"is_blocked"`
	IsDeleted    bool      `db:"is_deleted"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Role struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type Checklist struct {
	ID            uuid.UUID `db:"id"`
	Name          string    `db:"name"`
	Code          string    `db:"code"`
	AllowedRoleID uuid.UUID `db:"allowed_role_id"`
	CreatedAt     time.Time `db:"created_at"`
}

type Variety struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	Priority    string    `db:"priority"`
	ImageURL    *string   `db:"image_url"`
	CreatedAt   time.Time `db:"created_at"`
}

type Phenophase struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	ImageURL    *string   `db:"image_url"`
	OrderIndex  int       `db:"order_index"`
	CreatedAt   time.Time `db:"created_at"`
}

type Report struct {
	ID              uuid.UUID       `db:"id"`
	UserID          uuid.UUID       `db:"user_id"`
	ChecklistID     uuid.UUID       `db:"checklist_id"`
	VarietyID       *uuid.UUID      `db:"variety_id"`
	PhenophaseID    *uuid.UUID      `db:"phenophase_id"`
	ReportDate      time.Time       `db:"report_date"`
	ResponsibleName string          `db:"responsible_name"`
	Metadata        json.RawMessage `db:"metadata"`
	CreatedAt       time.Time       `db:"created_at"`
}

type Question struct {
	ID          uuid.UUID `db:"id"`
	Text        string    `db:"text"`
	OrderIndex  int       `db:"order_index"`
	IsActive    bool      `db:"is_active"`
	ChecklistID uuid.UUID `db:"checklist_id"`
	Formula     *string   `db:"formula"`
	ImageURL    *string   `db:"image_url"`
	CreatedAt   time.Time `db:"created_at"`
}

type QuestionPhenophaseFormula struct {
	ID           uuid.UUID `db:"id"`
	QuestionID   uuid.UUID `db:"question_id"`
	PhenophaseID uuid.UUID `db:"phenophase_id"`
	Formula      string    `db:"formula"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Answer struct {
	ID         uuid.UUID `db:"id"`
	ReportID   uuid.UUID `db:"report_id"`
	QuestionID uuid.UUID `db:"question_id"`
	AnswerText string    `db:"answer_text"`
	ImageURL   *string   `db:"image_url"`
	Result     *string   `db:"result"`
	CreatedAt  time.Time `db:"created_at"`
}
