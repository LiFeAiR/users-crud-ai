package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateRole общий метод для создания роли
func (bh *BaseHandler) CreateRole(ctx context.Context, in *api_pb.RoleCreateRequest) (out *api_pb.Role, err error) {
	// Validate request
	if in.Name == "" || in.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "Name and code are required")
	}

	role := models.Role{
		Name:        in.Name,
		Code:        in.Code,
		Description: in.Description,
	}

	// Используем репозиторий для создания роли
	dbRole, err := bh.roleRepo.CreateRole(ctx, &role)
	if err != nil {
		log.Printf("Failed to create role, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to create role")
	}

	// Send response
	return &api_pb.Role{
		Id:          int32(dbRole.ID),
		Name:        dbRole.Name,
		Code:        dbRole.Code,
		Description: dbRole.Description,
	}, nil
}