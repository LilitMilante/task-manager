package service

import (
	"context"
	"time"

	"task-manager/internal/entity"

	"github.com/google/uuid"
)

type TaskRepository interface {
	AddTask(ctx context.Context, task entity.Task) error
	TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error)
	Tasks(ctx context.Context) ([]entity.Task, error)
	UpdateTask(ctx context.Context, updateTask entity.TaskUpdated) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) AddTask(ctx context.Context, task entity.Task) (entity.Task, error) {
	task.ID = uuid.New()
	t := time.Now().UTC()
	task.CreatedAt = t
	task.UpdatedAt = t

	err := s.repo.AddTask(ctx, task)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (s *TaskService) TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	return s.repo.TaskByID(ctx, id)
}

func (s *TaskService) Tasks(ctx context.Context) ([]entity.Task, error) {
	return s.repo.Tasks(ctx)
}

func (s *TaskService) UpdateTask(ctx context.Context, updateTask entity.TaskUpdated) error {
	_, err := s.repo.TaskByID(ctx, updateTask.ID)
	if err != nil {
		return err
	}

	updateTask.UpdatedAt = time.Now().UTC()

	return s.repo.UpdateTask(ctx, updateTask)
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.TaskByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteTask(ctx, id)
}
