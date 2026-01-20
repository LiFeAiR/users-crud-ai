package handlers

import (
	"context"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/utils"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login общий метод для авторизации пользователя
func (bh *BaseHandler) Login(ctx context.Context, in *api_pb.LoginRequest) (out *api_pb.LoginResponse, err error) {
	// Проверяем входные данные
	if in.Email == "" || in.Password == "" {
		log.Println("Failed to validate login credentials")
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Ищем пользователя по email
	user, err := bh.userRepo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		log.Printf("Failed to get user by email, err:%v\n", err)
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	if user == nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	// Проверяем пароль
	isValid, err := bh.userRepo.CheckPassword(ctx, user.ID, in.Password)
	if err != nil {
		log.Printf("Failed to check password, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Authentication failed")
	}

	if !isValid {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	// Получаем права пользователя
	permissions, err := bh.userRepo.GetUserPermissions(ctx, user.ID)

	// Формируем ответ
	jwtPermissions := make([]string, 0, len(permissions))
	permissionsOut := make([]*api_pb.Permission, 0, len(permissions))
	for _, permission := range permissions {
		jwtPermissions = append(jwtPermissions, permission.Code)
		permissionsOut = append(permissionsOut, &api_pb.Permission{
			Id:          int32(permission.ID),
			Name:        permission.Name,
			Code:        permission.Code,
			Description: permission.Description,
		})
	}

	// Генерируем JWT токен
	token, err := utils.GenerateJWT(bh.secretKey, user.ID, user.Email, user.Name, jwtPermissions)
	if err != nil {
		log.Printf("Failed to generate JWT token, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Authentication failed")
	}

	// Формируем ответ
	var orgOut *api_pb.Organization
	if user.Organization != nil {
		org, err := bh.orgRepo.GetOrganizationByID(ctx, user.Organization.ID)
		if err == nil {
			user.Organization.Name = org.Name
		}

		orgOut = &api_pb.Organization{
			Id:   int32(user.Organization.ID),
			Name: user.Organization.Name,
		}
	}

	return &api_pb.LoginResponse{
		Token: token,
		User: &api_pb.User{
			Id:           int32(user.ID),
			Name:         user.Name,
			Email:        user.Email,
			Organization: orgOut,
			Permissions:  permissionsOut,
		},
	}, nil
}
