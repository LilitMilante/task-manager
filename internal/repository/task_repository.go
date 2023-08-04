package repository

import (
	"context"
	"database/sql"

	"task-manager/internal/api/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddTask(ctx context.Context, task entity.Task) error {
	q := `INSERT INTO tasks (id, name, description, created_at) VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, q, task.ID, task.Name, task.Description, task.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
