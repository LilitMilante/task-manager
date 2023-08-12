package entity

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	Name        string
	Description string
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskUpdated struct {
	ID          uuid.UUID
	Name        string
	Description string
	IsCompleted bool
	UpdatedAt   time.Time
}
