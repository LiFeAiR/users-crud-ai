package handlers

import (
	"context"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPermissions получает список прав
func (bh *BaseHandler) GetPermissions(ctx context.Context, in *api_pb.ListRequest) (out *api_pb.PermissionsResponse, err error) {
	// Проверяем входные данные
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

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

	// Используем репозиторий для получения списка прав
	permissions, err := bh.permRepo.GetPermissions(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get permissions")
	}

	// Формируем ответ
	var permissionsOut []*api_pb.Permission
	for _, permission := range permissions {
		permissionsOut = append(permissionsOut, &api_pb.Permission{
			Id:          int32(permission.ID),
			Name:        permission.Name,
			Code:        permission.Code,
			Description: permission.Description,
		})
	}

	// Возвращаем ответ
	return &api_pb.PermissionsResponse{
		Data: permissionsOut,
	}, nil
}
