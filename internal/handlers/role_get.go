package handlers

import (
	"context"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRole общий метод для получения роли по ID
func (bh *BaseHandler) GetRole(ctx context.Context, in *api_pb.Id) (out *api_pb.Role, err error) {
	// Проверяем входные данные
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Используем репозиторий для получения роли с правами
	role, err := bh.roleRepo.GetRoleWithPermissions(ctx, int(in.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "Role not found")
	}

	// Формируем ответ
	var permissionsOut []*api_pb.Permission
	for _, permission := range role.Permissions {
		permissionsOut = append(permissionsOut, &api_pb.Permission{
			Id:          int32(permission.ID),
			Name:        permission.Name,
			Code:        permission.Code,
			Description: permission.Description,
		})
	}

	// Возвращаем ответ
	return &api_pb.Role{
		Id:          int32(role.ID),
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Permissions: permissionsOut,
	}, nil
}
