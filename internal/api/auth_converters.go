package api

import (
	userv1 "task-manager/gen/proto/auth/v1"
	"task-manager/internal/entity"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserFromAPI(user *userv1.SignUpRequest) entity.User {
	return entity.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

func UserToAPI(user entity.User) *userv1.User {
	return &userv1.User{
		Id:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
