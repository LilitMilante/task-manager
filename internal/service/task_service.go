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
	Tasks(ctx context.Context) ([]entity.Task, error)
	UpdateTask(ctx context.Context, id uuid.UUID, updateTask entity.TaskUpdated) error
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
	t := time.Now().UTC()
	task.CreatedAt = t
	task.EditedAt = t

	err := s.repo.AddTask(ctx, task)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (s *Service) TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	return s.repo.TaskByID(ctx, id)
}

func (s *Service) Tasks(ctx context.Context) ([]entity.Task, error) {
	return s.repo.Tasks(ctx)
}

func (s *Service) UpdateTask(ctx context.Context, id uuid.UUID, updateTask entity.TaskUpdated) error {
	updateTask.EditedAt = time.Now().UTC()

	return s.repo.UpdateTask(ctx, id, updateTask)
}
