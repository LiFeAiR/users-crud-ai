package handlers

import (
	"context"
	"log"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddTariffRoles добавляет роли к тарифу
func (bh *BaseHandler) AddTariffRoles(
	ctx context.Context,
	in *api_pb.TariffRolesRequest,
) (out *api_pb.RolesResponse, err error) {
	// Проверяем входные данные
	if in == nil || in.Id == 0 || len(in.RoleIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Добавляем роли к тарифу
	if err := bh.tariffRepo.AddTariffRoles(ctx, int(in.Id), convertInt32SliceToInt(in.RoleIds)); err != nil {
		log.Printf("add tariff roles failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to add tariff roles")
	}

	// Получаем обновленный тариф с ролями
	roles, err := bh.tariffRepo.GetTariffWithRoles(ctx, int(in.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "Tariff not found")
	}

	// Формируем ответ
	var rolesOut []*api_pb.Role
	for _, role := range roles.Roles {
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