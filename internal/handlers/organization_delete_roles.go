package handlers

import (
	"context"
	"log"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteOrganizationRoles удаляет роли из организации
func (bh *BaseHandler) DeleteOrganizationRoles(
	ctx context.Context,
	in *api_pb.OrganizationRolesRequest,
) (out *api_pb.RolesResponse, err error) {
	// Проверяем входные данные
	if in == nil || in.Id == 0 || len(in.RoleIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Удаляем роли из организации
	if err := bh.orgRepo.DeleteOrganizationRoles(ctx, int(in.Id), convertInt32SliceToInt(in.RoleIds)); err != nil {
		log.Printf("delete organization roles failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to delete organization roles")
	}

	// Получаем обновленную организацию с ролями
	roles, err := bh.orgRepo.GetOrganizationRoles(ctx, int(in.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "Organization not found")
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
