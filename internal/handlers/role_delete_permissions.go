package handlers

import (
	"context"
	"log"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteRolePermissions удаляет права из роли
func (bh *BaseHandler) DeleteRolePermissions(
	ctx context.Context,
	in *api_pb.RolePermissionsRequest,
) (out *api_pb.RolePermissionsResponse, err error) {
	// Проверяем входные данные
	if in == nil || in.Id == 0 || len(in.PermissionIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Удаляем права из роли
	if err := bh.roleRepo.DeleteRolePermissions(ctx, int(in.Id), convertInt32SliceToInt(in.PermissionIds)); err != nil {
		log.Printf("delete role permissions failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to delete role permissions")
	}

	// Получаем обновленную роль с правами
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
	return &api_pb.RolePermissionsResponse{
		Data: permissionsOut,
	}, nil
}
