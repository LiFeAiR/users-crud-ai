package handlers

import (
	"context"

	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeletePermission общий метод для удаления права
func (bh *BaseHandler) DeletePermission(ctx context.Context, req *api_pb.Id) (*api_pb.Empty, error) {
	// Проверяем входные данные
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Используем репозиторий для удаления права
	if err := bh.permRepo.DeletePermission(ctx, int(req.GetId())); err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete permission")
	}

	// Возвращаем пустой ответ
	return &api_pb.Empty{}, nil
}
