package api

import (
	todolistv1 "task-manager/gen/proto/task/v1"
	"task-manager/internal/entity"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TaskFromAPI(request *todolistv1.AddTaskRequest) entity.Task {
	return entity.Task{
		Name:        request.Name,
		Description: request.Description,
	}
}

func TaskToAPI(task entity.Task) *todolistv1.Task {
	return &todolistv1.Task{
		Id:          task.ID.String(),
		UserId:      task.UserID.String(),
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

func TasksToAPI(tasks []entity.Task) []*todolistv1.Task {
	resp := make([]*todolistv1.Task, 0, len(tasks))
	for _, v := range tasks {
		resp = append(resp, TaskToAPI(v))
	}

	return resp
}
