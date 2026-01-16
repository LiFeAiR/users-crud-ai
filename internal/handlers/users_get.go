package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersResponse struct {
	Data []*models.User `json:"data"`
}

// GetUsersHandler обработчик для получения списка пользователей
func (h *BaseHandler) GetUsers(
	ctx context.Context,
	in *api_pb.ListRequest,
) (out *api_pb.UsersResponse, err error) {
	// Устанавливаем значения по умолчанию
	limit := 10
	offset := 0

	// Парсим limit
	if in.Limit > 0 && in.Limit < 100 {
		limit = int(in.GetLimit())
	}

	// Парсим offset
	if in.Offset > 0 {
		offset = int(in.GetOffset())
	}

	// Получаем список пользователей из репозитория
	users, err := h.userRepo.GetUsers(ctx, limit, offset)
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		return nil, status.Error(codes.Internal, "Failed to get users")
	}

	// Отправляем ответ клиенту
	data := make([]*api_pb.User, len(users))
	for i, user := range users {
		var orgOut *api_pb.Organization
		if user.Organization != nil {
			orgOut = &api_pb.Organization{
				Id:   int32(user.Organization.ID),
				Name: user.Organization.Name,
			}
		}
		data[i] = &api_pb.User{
			Id:           int32(user.ID),
			Name:         user.Name,
			Email:        user.Email,
			Organization: orgOut,
		}
	}

	return &api_pb.UsersResponse{
		Data: data,
	}, nil
}
