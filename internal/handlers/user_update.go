package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser общий метод для обновления пользователя
func (bh *BaseHandler) UpdateUser(ctx context.Context, in *api_pb.UserUpdateRequest) (out *api_pb.User, err error) {
	// Проверяем входные данные
	if in == nil || in.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	var org *models.Organization
	if in.GetOrganizationId() != 0 {
		org = &models.Organization{ID: int(in.GetOrganizationId())}
	}

	// Преобразуем запрос в модель
	user := models.User{
		ID:           int(in.Id),
		Name:         in.Name,
		Email:        in.Email,
		Password:     in.Password,
		Organization: org,
	}

	// Используем репозиторий для обновления пользователя
	if err := bh.userRepo.UpdateUser(ctx, &user); err != nil {
		log.Printf("update user failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to update user")
	}

	var orgOut *api_pb.Organization
	if user.Organization != nil {
		org, _ := bh.orgRepo.GetOrganizationByID(ctx, user.Organization.ID)
		if org != nil {
			orgOut = &api_pb.Organization{
				Id:   int32(user.Organization.ID),
				Name: org.Name,
			}
		}
	}

	// Возвращаем ответ
	return &api_pb.User{
		Id:           int32(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		Organization: orgOut,
	}, nil
}
