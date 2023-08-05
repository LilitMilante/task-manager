package repository

import (
	"context"
	"database/sql"

	"task-manager/internal/api/entity"

	"github.com/google/uuid"
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

func (r *Repository) TaskByID(ctx context.Context, id uuid.UUID) (task entity.Task, err error) {
	q := `SELECT id, name, description, created_at FROM tasks WHERE id = $1`

	err = r.db.QueryRowContext(ctx, q, id).Scan(&task.ID, &task.Name, &task.Description, &task.CreatedAt)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *Repository) Tasks(ctx context.Context) (tasks []entity.Task, err error) {
	q := `SELECT id, name, description, created_at FROM tasks`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var task entity.Task

		err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.CreatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
