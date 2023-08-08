package entity

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	Name        string
	Description string
	Status      bool
	CreatedAt   time.Time
	EditedAt    time.Time
}

type TaskUpdated struct {
	Name        string
	Description string
	Status      bool
	EditedAt    time.Time
}
