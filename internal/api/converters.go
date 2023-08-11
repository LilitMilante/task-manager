package api

import (
	todolistv1 "task-manager/gen/proto/task/v1"
	"task-manager/internal/api/entity"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TaskFromAPI(request *todolistv1.AddTaskRequest) entity.Task {
	return entity.Task{
		Name:        request.Name,
		Description: request.Description,
	}
}

func TaskToAPI(task entity.Task) *todolistv1.AddTaskResponse {
	return &todolistv1.AddTaskResponse{
		Id:          task.ID.String(),
		Name:        task.Name,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}
}

func TaskIDToAPI(task entity.Task) *todolistv1.TaskByIDResponse {
	return &todolistv1.TaskByIDResponse{
		Id:          task.ID.String(),
		Name:        task.Name,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}
}

func UpdateTaskFromAPI(updateTask *todolistv1.UpdateTaskRequest) (entity.TaskUpdated, error) {
	id, err := uuid.Parse(updateTask.Id)
	if err != nil {
		return entity.TaskUpdated{}, err
	}

	task := entity.TaskUpdated{
		ID:          id,
		Name:        updateTask.Name,
		Description: updateTask.Description,
		IsCompleted: updateTask.IsCompleted,
	}

	return task, nil
}
