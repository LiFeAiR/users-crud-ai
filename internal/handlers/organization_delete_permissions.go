package handlers

import (
	"context"
	"log"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteOrganizationPermissions удаляет права из организации
func (bh *BaseHandler) DeleteOrganizationPermissions(
	ctx context.Context,
	in *api_pb.OrganizationPermissionsRequest,
) (out *api_pb.RolePermissionsResponse, err error) {
	// Проверяем входные данные
	if in == nil || in.Id == 0 || len(in.PermissionIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Удаляем права из организации
	if err := bh.orgRepo.DeleteOrganizationPermissions(ctx, int(in.Id), convertInt32SliceToInt(in.PermissionIds)); err != nil {
		log.Printf("delete organization permissions failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to delete organization permissions")
	}

	// Получаем обновленную организацию с правами
	permissions, err := bh.orgRepo.GetOrganizationPermissions(ctx, int(in.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "Organization not found")
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
