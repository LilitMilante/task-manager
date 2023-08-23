package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"task-manager/internal/entity"

	"github.com/google/uuid"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
	SignIn(ctx context.Context, email, password string) (uuid.UUID, error)
	CreateSession(ctx context.Context, session entity.Session) error
	UserByEmail(ctx context.Context, email string) error
	UserBySessionID(ctx context.Context, sessionID uuid.UUID) (entity.User, error)
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
	err := a.repo.UserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, err
	}

	user.ID = uuid.New()
	user.CreatedAt = time.Now().UTC()

	err = a.repo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	user.Password = ""

	return user, nil
}

func (a *AuthService) SignIn(ctx context.Context, email, password string) (uuid.UUID, error) {
	userID, err := a.repo.SignIn(ctx, email, password)
	if err != nil {
		return uuid.UUID{}, err
	}

	session := entity.Session{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}

	err = a.repo.CreateSession(ctx, session)
	if err != nil {
		return uuid.UUID{}, err
	}

	return session.ID, nil
}

func (a *AuthService) Auth(ctx context.Context, sessionID uuid.UUID) (entity.User, error) {
	return a.repo.UserBySessionID(ctx, sessionID)
}

func (a *AuthService) AuthorisedUser(ctx context.Context) (entity.User, error) {
	value := ctx.Value("user")

	user, ok := value.(entity.User)
	if !ok {
		return entity.User{}, ErrUnauthorized
	}

	return user, nil
}
