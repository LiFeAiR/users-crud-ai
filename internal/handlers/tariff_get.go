package handlers

import (
	"context"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetTariff общий метод для получения тарифа по ID
func (bh *BaseHandler) GetTariff(ctx context.Context, in *api_pb.Id) (out *api_pb.Tariff, err error) {
	// Проверяем входные данные
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Используем репозиторий для получения тарифа с ролями
	tariff, err := bh.tariffRepo.GetTariffWithRoles(ctx, int(in.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "Tariff not found")
	}

	// Формируем ответ
	var rolesOut []*api_pb.Role
	for _, role := range tariff.Roles {
		// TODO N+1 запросы
		rr, err := bh.roleRepo.GetRoleWithPermissions(ctx, role.ID)
		if err != nil {
			return nil, status.Error(codes.NotFound, "Role not found")
		}

		// Формируем ответ
		var permissionsOut []*api_pb.Permission
		for _, permission := range rr.Permissions {
			permissionsOut = append(permissionsOut, &api_pb.Permission{
				Id:          int32(permission.ID),
				Name:        permission.Name,
				Code:        permission.Code,
				Description: permission.Description,
			})
		}

		rolesOut = append(rolesOut, &api_pb.Role{
			Id:          int32(role.ID),
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Permissions: permissionsOut,
		})
	}

	// Возвращаем ответ
	return &api_pb.Tariff{
		Id:          int32(tariff.ID),
		Name:        tariff.Name,
		Description: tariff.Description,
		Price:       int32(tariff.Price),
		Roles:       rolesOut,
	}, nil
}
