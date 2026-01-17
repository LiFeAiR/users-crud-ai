package handlers

import (
	"context"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetTariffs общий метод для получения списка тарифов
func (bh *BaseHandler) GetTariffs(ctx context.Context, in *api_pb.ListRequest) (out *api_pb.TariffsResponse, err error) {
	// Устанавливаем значения по умолчанию
	limit := 10
	offset := 0

	// Парсим limit
	if in.Limit > 0 && in.Limit < 100 {
		limit = int(in.GetLimit())
	}

	// Парсим offset
	if in.Offset > 0 {
		offset = int(in.GetOffset())
	}

	// Используем репозиторий для получения списка тарифов
	tariffs, err := bh.tariffRepo.GetTariffs(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get tariffs")
	}

	// Формируем ответ
	var tariffsOut []*api_pb.Tariff
	for _, tariff := range tariffs {
		tariffsOut = append(tariffsOut, &api_pb.Tariff{
			Id:          int32(tariff.ID),
			Name:        tariff.Name,
			Description: tariff.Description,
			Price:       int32(tariff.Price),
		})
	}

	// Возвращаем ответ
	return &api_pb.TariffsResponse{
		Data: tariffsOut,
	}, nil
}
