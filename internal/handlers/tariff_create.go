package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateTariff общий метод для создания тарифа
func (bh *BaseHandler) CreateTariff(ctx context.Context, in *api_pb.TariffCreateRequest) (out *api_pb.Tariff, err error) {
	// Validate request
	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is required")
	}

	tariff := models.Tariff{
		Name:        in.Name,
		Description: in.Description,
		Price:       int(in.Price),
	}

	// Используем репозиторий для создания тарифа
	dbTariff, err := bh.tariffRepo.CreateTariff(ctx, &tariff)
	if err != nil {
		log.Printf("Failed to create tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to create tariff")
	}

	// Send response
	return &api_pb.Tariff{
		Id:          int32(dbTariff.ID),
		Name:        dbTariff.Name,
		Description: dbTariff.Description,
		Price:       int32(dbTariff.Price),
	}, nil
}