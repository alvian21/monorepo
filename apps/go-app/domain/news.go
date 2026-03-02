package domain

import (
	"time"

	"github.com/google/uuid"
)

type NewsStatus string

const (
	NewsStatusDraft     NewsStatus = "DRAFT"
	NewsStatusPublished NewsStatus = "PUBLISHED"
	NewsStatusDeleted   NewsStatus = "DELETED"
)

type News struct {
	ID        uuid.UUID    `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Status    NewsStatus   `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Topics    []Topic      `json:"topics,omitempty"`
}

type CreateNewsRequest struct {
	Title    string      `json:"title" validate:"required"`
	Content  string      `json:"content" validate:"required"`
	TopicIDs []uuid.UUID `json:"topic_ids"`
}

type UpdateNewsRequest struct {
	Title    string      `json:"title" validate:"required"`
	Content  string      `json:"content" validate:"required"`
	Status   NewsStatus  `json:"status" validate:"required,oneof=DRAFT PUBLISHED DELETED"`
	TopicIDs []uuid.UUID `json:"topic_ids"`
}

type NewsFilter struct {
	Status NewsStatus `json:"status" query:"status"`
	Search string     `json:"search" query:"search"`
}

