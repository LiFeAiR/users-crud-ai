package handlers

import (
	"context"
	"log"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteUserPermissions удаляет права у пользователя
func (bh *BaseHandler) DeleteUserPermissions(
	ctx context.Context,
	in *api_pb.UserPermissionsRequest,
) (out *api_pb.RolePermissionsResponse, err error) {
	// Проверяем входные данные
	if in == nil || in.UserId == 0 || len(in.PermissionIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Удаляем права у пользователя
	if err := bh.userRepo.DeleteUserPermissions(ctx, int(in.UserId), convertInt32SliceToInt(in.PermissionIds)); err != nil {
		log.Printf("delete user permissions failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to delete user permissions")
	}

	// Получаем обновленного пользователя с правами
	permissions, err := bh.userRepo.GetUserPermissions(ctx, int(in.UserId))
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
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
	return &api_pb.RolePermissionsResponse{
		Data: permissionsOut,
	}, nil
}