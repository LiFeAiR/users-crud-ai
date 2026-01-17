package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdatePermission общий метод для обновления права
func (bh *BaseHandler) UpdatePermission(ctx context.Context, in *api_pb.PermissionUpdateRequest) (out *api_pb.Permission, err error) {
	// Проверяем входные данные
	if in == nil || in.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Преобразуем запрос в модель
	permission := models.Permission{
		ID:          int(in.Id),
		Name:        in.Name,
		Code:        in.Code,
		Description: in.Description,
	}

	// Используем репозиторий для обновления права
	if err := bh.permRepo.UpdatePermission(ctx, &permission); err != nil {
		log.Printf("update permission failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to update permission")
	}

	// Возвращаем ответ
	return &api_pb.Permission{
		Id:          int32(permission.ID),
		Name:        permission.Name,
		Code:        permission.Code,
		Description: permission.Description,
	}, nil
}
