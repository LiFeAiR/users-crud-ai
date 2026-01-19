package handlers

import (
	"context"
	"log"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddUserTariff добавляет тариф пользователю
func (bh *BaseHandler) AddUserTariff(ctx context.Context, in *api_pb.UserTariffRequest) (out *api_pb.UserTariffResponse, err error) {
	// Проверяем входные данные
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}
	if in.TariffId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid tariff_id")
	}

	// Проверяем, существует ли пользователь
	_, err = bh.userRepo.GetUserByID(ctx, int(in.Id))
	if err != nil {
		log.Printf("Failed to get user, err:%v\n", err)
		return nil, status.Error(codes.NotFound, "User not found")
	}

	// Проверяем, существует ли тариф
	_, err = bh.tariffRepo.GetTariffByID(ctx, int(in.TariffId))
	if err != nil {
		log.Printf("Failed to get tariff, err:%v\n", err)
		return nil, status.Error(codes.NotFound, "Tariff not found")
	}

	// Устанавливаем тариф пользователю
	err = bh.userRepo.SetUserTariff(ctx, int(in.Id), &in.TariffId)
	if err != nil {
		log.Printf("Failed to set user tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to set user tariff")
	}

	// Получаем тариф для ответа
	tariff, err := bh.tariffRepo.GetTariffByID(ctx, int(in.TariffId))
	if err != nil {
		log.Printf("Failed to get tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to get tariff")
	}

	// Формируем ответ
	return &api_pb.UserTariffResponse{
		Tariff: &api_pb.Tariff{
			Id:          int32(tariff.ID),
			Name:        tariff.Name,
			Description: tariff.Description,
			Price:       int32(tariff.Price),
		},
	}, nil
}

// UpdateUserTariff обновляет тариф пользователя
func (bh *BaseHandler) UpdateUserTariff(ctx context.Context, in *api_pb.UserTariffRequest) (out *api_pb.UserTariffResponse, err error) {
	// Проверяем входные данные
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}
	if in.TariffId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid tariff_id")
	}

	// Проверяем, существует ли пользователь
	_, err = bh.userRepo.GetUserByID(ctx, int(in.Id))
	if err != nil {
		log.Printf("Failed to get user, err:%v\n", err)
		return nil, status.Error(codes.NotFound, "User not found")
	}

	// Проверяем, существует ли тариф
	_, err = bh.tariffRepo.GetTariffByID(ctx, int(in.TariffId))
	if err != nil {
		log.Printf("Failed to get tariff, err:%v\n", err)
		return nil, status.Error(codes.NotFound, "Tariff not found")
	}

	// Обновляем тариф пользователя
	err = bh.userRepo.SetUserTariff(ctx, int(in.Id), &in.TariffId)
	if err != nil {
		log.Printf("Failed to update user tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to update user tariff")
	}

	// Получаем тариф для ответа
	tariff, err := bh.tariffRepo.GetTariffByID(ctx, int(in.TariffId))
	if err != nil {
		log.Printf("Failed to get tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to get tariff")
	}

	// Формируем ответ
	return &api_pb.UserTariffResponse{
		Tariff: &api_pb.Tariff{
			Id:          int32(tariff.ID),
			Name:        tariff.Name,
			Description: tariff.Description,
			Price:       int32(tariff.Price),
		},
	}, nil
}

// DeleteUserTariff удаляет тариф пользователя
func (bh *BaseHandler) DeleteUserTariff(ctx context.Context, in *api_pb.UserTariffRequest) (out *api_pb.Empty, err error) {
	// Проверяем входные данные
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	// Проверяем, существует ли пользователь
	_, err = bh.userRepo.GetUserByID(ctx, int(in.Id))
	if err != nil {
		log.Printf("Failed to get user, err:%v\n", err)
		return nil, status.Error(codes.NotFound, "User not found")
	}

	// Удаляем тариф пользователя (устанавливаем nil)
	err = bh.userRepo.SetUserTariff(ctx, int(in.Id), nil)
	if err != nil {
		log.Printf("Failed to delete user tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to delete user tariff")
	}

	return &api_pb.Empty{}, nil
}
