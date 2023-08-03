package service

import (
	"context"
	"time"

	"task-manager/internal/api/entity"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AddTask(ctx context.Context, task entity.Task) (entity.Task, error) {
	return entity.Task{
		ID:          1,
		Name:        task.Name,
		Description: task.Description,
		CreatedAt:   time.Now(),
	}, nil
}
