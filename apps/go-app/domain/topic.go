package domain

import (
	"time"

	"github.com/google/uuid"
)

type Topic struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTopicRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type UpdateTopicRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type TopicFilter struct {
	Search string `json:"search" query:"search"`
}

