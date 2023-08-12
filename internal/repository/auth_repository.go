package repository

import (
	"context"

	"task-manager/internal/entity"

	"github.com/google/uuid"
)

func (r *Repository) CreateUser(ctx context.Context, user entity.User) error {
	const q = "INSERT INTO users (id, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5)"

	_, err := r.db.ExecContext(ctx, q, user.ID, user.Name, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UserByEmail(ctx context.Context, email string) error {
	var id uuid.UUID
	const q = "SELECT id from users WHERE email = $1"

	return r.db.QueryRowContext(ctx, q, email).Scan(&id)
}

func (r *Repository) SignIn(ctx context.Context, email, password string) (uuid.UUID, error) {
	var id uuid.UUID
	const q = "SELECT id from users WHERE email = $1 AND password = $2"

	err := r.db.QueryRowContext(ctx, q, email, password).Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (r *Repository) CreateSession(ctx context.Context, session entity.Session) (uuid.UUID, error) {
	const q = "INSERT INTO sessions (session_id, user_id, created_at) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, q, session.ID, session.UserID, session.CreatedAt)
	if err != nil {
		return uuid.UUID{}, err
	}

	return session.ID, nil
}
