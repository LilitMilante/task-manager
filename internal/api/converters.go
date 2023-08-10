package api

import (
	todolistv1 "task-manager/gen/proto/task/v1"
	"task-manager/internal/api/entity"

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
