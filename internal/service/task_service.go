package service

import (
	"context"
	"time"

	"task-manager/internal/api/entity"

	"github.com/google/uuid"
)

type Repository interface {
	AddTask(ctx context.Context, task entity.Task) error
	TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddTask(ctx context.Context, task entity.Task) (entity.Task, error) {
	task.ID = uuid.New()
	task.CreatedAt = time.Now()

	err := s.repo.AddTask(ctx, task)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (s *Service) TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	return s.repo.TaskByID(ctx, id)
}
