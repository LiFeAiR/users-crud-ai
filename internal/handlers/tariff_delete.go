package handlers

import (
	"context"
	"log"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteTariff общий метод для удаления тарифа
func (bh *BaseHandler) DeleteTariff(ctx context.Context, in *api_pb.Id) (out *api_pb.Empty, err error) {
	// Проверяем входные данные
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Используем репозиторий для удаления тарифа
	err = bh.tariffRepo.DeleteTariff(ctx, int(in.Id))
	if err != nil {
		log.Printf("Failed to delete tariff, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to delete tariff")
	}

	// Возвращаем пустой ответ
	return &api_pb.Empty{}, nil
}