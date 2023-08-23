package repository

import (
	"context"
	"database/sql"
	"errors"

	"task-manager/internal/entity"
	"task-manager/internal/service"

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
	const q = `INSERT INTO tasks (id, user_id, name, description, is_completed, created_at, updated_at)
			   VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, q, task.ID, task.UserID, task.Name, task.Description, task.IsCompleted, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) TaskByID(ctx context.Context, id uuid.UUID) (task entity.Task, err error) {
	const q = `SELECT id, user_id, name, description, is_completed, created_at, updated_at
			   FROM tasks WHERE id = $1`

	err = r.db.QueryRowContext(ctx, q, id).Scan(&task.ID, &task.UserID, &task.Name, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Task{}, service.ErrNotFound
		}

		return entity.Task{}, err
	}

	return task, nil
}

func (r *Repository) Tasks(ctx context.Context, userID uuid.UUID) (tasks []entity.Task, err error) {
	const q = `SELECT id, user_id, name, description, is_completed, created_at, updated_at
		  	   FROM tasks
		  	   WHERE user_id = $1
		  	   ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, q, userID)
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

		err := rows.Scan(&task.ID, &task.UserID, &task.Name, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *Repository) UpdateTask(ctx context.Context, updateTask entity.TaskUpdated) error {
	const q = `UPDATE tasks SET name = $1, description = $2, is_completed = $3, updated_at = $4 WHERE id = $5`

	_, err := r.db.ExecContext(ctx, q, updateTask.Name, updateTask.Description, updateTask.IsCompleted, updateTask.UpdatedAt, updateTask.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM tasks WHERE id = $1`

	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}
