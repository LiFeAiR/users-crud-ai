package handlers

import (
	"context"
	"log"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddOrganizationTariff добавляет тариф организации
func (bh *BaseHandler) AddOrganizationTariff(ctx context.Context, in *api_pb.OrganizationTariffRequest) (out *api_pb.OrganizationTariffResponse, err error) {
	// Проверяем входные данные
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}
	if in.TariffId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid tariff_id")
	}

	// Проверяем, существует ли организация
	_, err = bh.orgRepo.GetOrganizationByID(ctx, int(in.Id))
	if err != nil {
		log.Printf("Failed to get organization, err:%v\n", err)
		return nil, status.Error(codes.NotFound, "Organization not found")
	}

	// Проверяем, существует ли тариф
	_, err = bh.tariffRepo.GetTariffByID(ctx, int(in.TariffId))
	if err != nil {
		log.Printf("Failed to get tariff, err:%v\n", err)
		return nil, status.Error(codes.NotFound, "Tariff not found")
	}

	// Устанавливаем тариф организации
	err = bh.orgRepo.SetOrganizationTariff(ctx, int(in.Id), &in.TariffId)
	if err != nil {
		log.Printf("Failed to set organization tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to set organization tariff")
	}

	// Получаем тариф для ответа
	tariff, err := bh.tariffRepo.GetTariffByID(ctx, int(in.TariffId))
	if err != nil {
		log.Printf("Failed to get tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to get tariff")
	}

	// Формируем ответ
	return &api_pb.OrganizationTariffResponse{
		Tariff: &api_pb.Tariff{
			Id:          int32(tariff.ID),
			Name:        tariff.Name,
			Description: tariff.Description,
			Price:       int32(tariff.Price),
		},
	}, nil
}

// UpdateOrganizationTariff обновляет тариф организации
func (bh *BaseHandler) UpdateOrganizationTariff(ctx context.Context, in *api_pb.OrganizationTariffRequest) (out *api_pb.OrganizationTariffResponse, err error) {
	// Проверяем входные данные
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}
	if in.TariffId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid tariff_id")
	}

	// Проверяем, существует ли организация
	_, err = bh.orgRepo.GetOrganizationByID(ctx, int(in.Id))
	if err != nil {
		log.Printf("Failed to get organization, err:%v\n", err)
		return nil, status.Error(codes.NotFound, "Organization not found")
	}

	// Проверяем, существует ли тариф
	_, err = bh.tariffRepo.GetTariffByID(ctx, int(in.TariffId))
	if err != nil {
		log.Printf("Failed to get tariff, err:%v\n", err)
		return nil, status.Error(codes.NotFound, "Tariff not found")
	}

	// Обновляем тариф организации
	err = bh.orgRepo.SetOrganizationTariff(ctx, int(in.Id), &in.TariffId)
	if err != nil {
		log.Printf("Failed to update organization tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to update organization tariff")
	}

	// Получаем тариф для ответа
	tariff, err := bh.tariffRepo.GetTariffByID(ctx, int(in.TariffId))
	if err != nil {
		log.Printf("Failed to get tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to get tariff")
	}

	// Формируем ответ
	return &api_pb.OrganizationTariffResponse{
		Tariff: &api_pb.Tariff{
			Id:          int32(tariff.ID),
			Name:        tariff.Name,
			Description: tariff.Description,
			Price:       int32(tariff.Price),
		},
	}, nil
}

// DeleteOrganizationTariff удаляет тариф организации
func (bh *BaseHandler) DeleteOrganizationTariff(ctx context.Context, in *api_pb.OrganizationTariffRequest) (out *api_pb.Empty, err error) {
	// Проверяем входные данные
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	// Проверяем, существует ли организация
	_, err = bh.orgRepo.GetOrganizationByID(ctx, int(in.Id))
	if err != nil {
		log.Printf("Failed to get organization, err:%v\n", err)
		return nil, status.Error(codes.NotFound, "Organization not found")
	}

	// Удаляем тариф организации (устанавливаем nil)
	err = bh.orgRepo.SetOrganizationTariff(ctx, int(in.Id), nil)
	if err != nil {
		log.Printf("Failed to delete organization tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to delete organization tariff")
	}

	return &api_pb.Empty{}, nil
}
