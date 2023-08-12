package repository

import (
	"context"

	"task-manager/internal/entity"
)

func (r *Repository) CreateUser(ctx context.Context, user entity.User) error {
	const q = "INSERT INTO users (id, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5)"

	_, err := r.db.ExecContext(ctx, q, user.ID, user.Name, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
