package service

import (
	"context"
	"fmt"
	"time"

	"task-manager/internal/entity"

	"github.com/google/uuid"
)

type TaskRepository interface {
	AddTask(ctx context.Context, task entity.Task) error
	TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error)
	Tasks(ctx context.Context, userID uuid.UUID) ([]entity.Task, error)
	UpdateTask(ctx context.Context, updateTask entity.TaskUpdated) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
}

type Auth interface {
	AuthorisedUser(ctx context.Context) (entity.User, error)
}

type TaskService struct {
	repo TaskRepository
	auth Auth
}

func NewTaskService(repo TaskRepository, auth Auth) *TaskService {
	return &TaskService{
		repo: repo,
		auth: auth,
	}
}

func (s *TaskService) AddTask(ctx context.Context, task entity.Task) (entity.Task, error) {
	user, err := s.auth.AuthorisedUser(ctx)
	if err != nil {
		return entity.Task{}, fmt.Errorf("can`t to create creat task: %w", err)
	}

	task.ID = uuid.New()
	task.UserID = user.ID
	t := time.Now().UTC()
	task.CreatedAt = t
	task.UpdatedAt = t

	err = s.repo.AddTask(ctx, task)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (s *TaskService) TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error) {
	user, err := s.auth.AuthorisedUser(ctx)
	if err != nil {
		return entity.Task{}, err
	}

	task, err := s.repo.TaskByID(ctx, id)
	if err != nil {
		return entity.Task{}, err
	}

	if task.UserID != user.ID {
		return entity.Task{}, ErrAccessDenied
	}

	return task, nil
}

func (s *TaskService) Tasks(ctx context.Context) ([]entity.Task, error) {
	user, err := s.auth.AuthorisedUser(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.Tasks(ctx, user.ID)
}

func (s *TaskService) UpdateTask(ctx context.Context, updateTask entity.TaskUpdated) error {
	_, err := s.TaskByID(ctx, updateTask.ID)
	if err != nil {
		return err
	}

	updateTask.UpdatedAt = time.Now().UTC()

	return s.repo.UpdateTask(ctx, updateTask)
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	_, err := s.TaskByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteTask(ctx, id)
}
