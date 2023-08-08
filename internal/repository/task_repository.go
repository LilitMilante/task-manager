package repository

import (
	"context"
	"database/sql"

	"task-manager/internal/api/entity"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository struct {
	l  *zap.SugaredLogger
	db *sql.DB
}

func NewRepository(l *zap.SugaredLogger, db *sql.DB) *Repository {
	return &Repository{
		l:  l,
		db: db,
	}
}

func (r *Repository) AddTask(ctx context.Context, task entity.Task) error {
	q := `INSERT INTO tasks (id, name, description, is_completed, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, q, task.ID, task.Name, task.Description, task.IsCompleted, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) TaskByID(ctx context.Context, id uuid.UUID) (task entity.Task, err error) {
	q := `SELECT id, name, description, is_completed, created_at, updated_at FROM tasks WHERE id = $1`

	err = r.db.QueryRowContext(ctx, q, id).Scan(&task.ID, &task.Name, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *Repository) Tasks(ctx context.Context) (tasks []entity.Task, err error) {
	q := `SELECT id, name, description, is_completed, created_at, updated_at FROM tasks`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.l.Warnw("rows close", zap.Error(err))
		}
	}(rows)

	for rows.Next() {
		var task entity.Task

		err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *Repository) UpdateTask(ctx context.Context, id uuid.UUID, updateTask entity.TaskUpdated) error {
	q := `UPDATE tasks SET name = $1, description = $2, is_completed = $3, updated_at = $4 WHERE id = $5`

	_, err := r.db.ExecContext(ctx, q, updateTask.Name, updateTask.Description, updateTask.IsCompleted, updateTask.UpdatedAt, id)
	if err != nil {
		return err
	}

	return nil
}
