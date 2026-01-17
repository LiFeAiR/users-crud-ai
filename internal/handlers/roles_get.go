package handlers

import (
	"context"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRoles получает список ролей
func (bh *BaseHandler) GetRoles(ctx context.Context, in *api_pb.ListRequest) (out *api_pb.RolesResponse, err error) {
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

	// Используем репозиторий для получения списка ролей
	roles, err := bh.roleRepo.GetRoles(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get roles")
	}

	// Формируем ответ
	var rolesOut []*api_pb.Role
	for _, role := range roles {
		rolesOut = append(rolesOut, &api_pb.Role{
			Id:          int32(role.ID),
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
		})
	}

	// Возвращаем ответ
	return &api_pb.RolesResponse{
		Data: rolesOut,
	}, nil
}
