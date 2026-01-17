package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateTariff общий метод для обновления тарифа
func (bh *BaseHandler) UpdateTariff(ctx context.Context, in *api_pb.TariffUpdateRequest) (out *api_pb.Tariff, err error) {
	// Validate request
	if in.Id == 0 || in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Id and Name are required")
	}

	tariff := models.Tariff{
		ID:          int(in.Id),
		Name:        in.Name,
		Description: in.Description,
		Price:       int(in.Price),
	}

	// Используем репозиторий для обновления тарифа
	err = bh.tariffRepo.UpdateTariff(ctx, &tariff)
	if err != nil {
		log.Printf("Failed to update tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to update tariff")
	}

	// Получаем обновленный тариф
	updatedTariff, err := bh.tariffRepo.GetTariffByID(ctx, int(in.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "Tariff not found")
	}

	// Send response
	return &api_pb.Tariff{
		Id:          int32(updatedTariff.ID),
		Name:        updatedTariff.Name,
		Description: updatedTariff.Description,
		Price:       int32(updatedTariff.Price),
	}, nil
}