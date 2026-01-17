package handlers

import (
	"context"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPermission общий метод для получения права по ID
func (bh *BaseHandler) GetPermission(ctx context.Context, in *api_pb.Id) (out *api_pb.Permission, err error) {
	// Проверяем входные данные
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Используем репозиторий для получения права
	permission, err := bh.permRepo.GetPermissionByID(ctx, int(in.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "Permission not found")
	}

	// Возвращаем ответ
	return &api_pb.Permission{
		Id:          int32(permission.ID),
		Name:        permission.Name,
		Code:        permission.Code,
		Description: permission.Description,
	}, nil
}
