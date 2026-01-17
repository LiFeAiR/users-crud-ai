package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreatePermission общий метод для создания права
func (bh *BaseHandler) CreatePermission(ctx context.Context, in *api_pb.PermissionCreateRequest) (out *api_pb.Permission, err error) {
	// Validate request
	if in.Name == "" || in.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "Name and code are required")
	}

	permission := models.Permission{
		Name:        in.Name,
		Code:        in.Code,
		Description: in.Description,
	}

	// Используем репозиторий для создания права
	dbPermission, err := bh.permRepo.CreatePermission(ctx, &permission)
	if err != nil {
		log.Printf("Failed to create permission, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to create permission")
	}

	// Send response
	return &api_pb.Permission{
		Id:          int32(dbPermission.ID),
		Name:        dbPermission.Name,
		Code:        dbPermission.Code,
		Description: dbPermission.Description,
	}, nil
}
