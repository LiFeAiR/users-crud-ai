package handlers

import (
	"context"

	"github.com/LiFeAiR/crud-ai/internal/utils"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetOrganization получает организацию по ID
func (bh *BaseHandler) GetOrganization(ctx context.Context, in *api_pb.Id) (out *api_pb.Organization, err error) {
	// Проверяем входные данные
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Используем репозиторий для получения организации
	org, err := bh.orgRepo.GetOrganizationByID(ctx, int(in.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "Organization not found")
	}

	// Получаем обновленную организацию с правами
	permissions, err := bh.orgRepo.GetOrganizationPermissions(ctx, int(in.Id))

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
	return &api_pb.Organization{
		Id:          int32(org.ID),
		Name:        org.Name,
		TariffId:    int32(utils.FromPtr(org.TariffID)),
		Permissions: permissionsOut,
	}, nil
}
