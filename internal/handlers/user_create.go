package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser общий метод для создания пользователя
func (bh *BaseHandler) CreateUser(ctx context.Context, in *api_pb.UserCreateRequest) (out *api_pb.User, err error) {
	var org *models.Organization
	if in.GetOrganizationId() != 0 {
		org = &models.Organization{ID: int(in.GetOrganizationId())}
	}

	user := models.User{
		Name:         in.Name,
		Email:        in.Email,
		Password:     in.Password,
		Organization: org,
	}

	// Validate request
	if false {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Используем репозиторий для создания пользователя
	dbUser, err := bh.userRepo.CreateUser(ctx, &user)
	if err != nil {
		log.Printf("Failed to create user, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	var orgOut *api_pb.Organization
	if user.Organization != nil {
		org, _ := bh.orgRepo.GetOrganizationByID(ctx, user.Organization.ID)
		if org != nil {
			orgOut = &api_pb.Organization{
				Id:   int32(user.Organization.ID),
				Name: org.Name,
			}
		}
	}

	// Send response
	return &api_pb.User{
		Id:           int32(dbUser.ID),
		Name:         user.Name,
		Email:        user.Email,
		Organization: orgOut,
	}, nil
}
