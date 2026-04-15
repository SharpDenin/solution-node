package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	FullName     string    `db:"full_name"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
}

type Place struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type Report struct {
	ID              uuid.UUID `db:"id"`
	UserID          uuid.UUID `db:"user_id"`
	PlaceID         uuid.UUID `db:"place_id"`
	ReportDate      time.Time `db:"report_date"`
	ResponsibleName string    `db:"responsible_name"`
	CreatedAt       time.Time `db:"created_at"`
}

type Question struct {
	ID         uuid.UUID `db:"id"`
	Text       string    `db:"text"`
	OrderIndex int       `db:"order_index"`
	IsActive   bool      `db:"is_active"`
	CreatedAt  time.Time `db:"created_at"`
}

type Answer struct {
	ID         uuid.UUID `db:"id"`
	ReportID   uuid.UUID `db:"report_id"`
	QuestionID uuid.UUID `db:"question_id"`
	AnswerText string    `db:"answer_text"`
	ImageURL   *string   `db:"image_url"`
	CreatedAt  time.Time `db:"created_at"`
}
