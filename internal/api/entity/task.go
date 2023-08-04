package entity

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
}
