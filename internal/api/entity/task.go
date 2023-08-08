package entity

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskUpdated struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	UpdatedAt   time.Time `json:"updated_at"`
}
