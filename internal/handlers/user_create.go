package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/LiFeAiR/crud-ai/internal/utils"
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

	// Validate request
	// TODO валидация на сложный пароль
	if in.Password == "" || len(in.Password) < 5 {
		log.Println("Failed to validate password")
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Хешируем пароль перед сохранением
	hash, err := utils.HashPassword(in.Password)
	if err != nil {
		log.Printf("Failed to hash password, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	user := models.User{
		Name:         in.Name,
		Email:        in.Email,
		PasswordHash: hash,
		Organization: org,
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
