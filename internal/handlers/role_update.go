package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateRole общий метод для обновления роли
func (bh *BaseHandler) UpdateRole(ctx context.Context, in *api_pb.RoleUpdateRequest) (out *api_pb.Role, err error) {
	// Проверяем входные данные
	if in == nil || in.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Преобразуем запрос в модель
	role := models.Role{
		ID:          int(in.Id),
		Name:        in.Name,
		Code:        in.Code,
		Description: in.Description,
	}

	// Используем репозиторий для обновления роли
	if err := bh.roleRepo.UpdateRole(ctx, &role); err != nil {
		log.Printf("update role failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to update role")
	}

	// Возвращаем ответ
	return &api_pb.Role{
		Id:          int32(role.ID),
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
	}, nil
}