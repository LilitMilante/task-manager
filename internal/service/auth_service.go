package service

import (
	"context"
	"time"

	"task-manager/internal/entity"

	"github.com/google/uuid"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) SignUp(ctx context.Context, user entity.User) (entity.User, error) {
	//проверить сузесвует ли

	user.ID = uuid.New()
	user.CreatedAt = time.Now().UTC()

	err := a.repo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	user.Password = ""

	return user, nil
}
