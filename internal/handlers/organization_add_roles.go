package handlers

import (
	"context"
	"log"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddOrganizationRoles добавляет роли к организации
func (bh *BaseHandler) AddOrganizationRoles(
	ctx context.Context,
	in *api_pb.OrganizationRolesRequest,
) (out *api_pb.RolesResponse, err error) {
	// Проверяем входные данные
	if in == nil || in.Id == 0 || len(in.RoleIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Добавляем роли к организации
	if err := bh.orgRepo.AddOrganizationRoles(ctx, int(in.Id), convertInt32SliceToInt(in.RoleIds)); err != nil {
		log.Printf("add organization roles failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to add organization roles")
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
