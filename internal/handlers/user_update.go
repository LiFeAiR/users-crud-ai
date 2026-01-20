package handlers

import (
	"context"
	"errors"
	"log"

	"github.com/LiFeAiR/crud-ai/internal/models"
	"github.com/LiFeAiR/crud-ai/internal/server/middleware/auth"
	"github.com/LiFeAiR/crud-ai/internal/utils"
	api_pb "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser общий метод для обновления пользователя
func (bh *BaseHandler) UpdateUser(ctx context.Context, in *api_pb.UserUpdateRequest) (out *api_pb.User, err error) {
	// Проверяем входные данные
	if in == nil || in.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	err = checkPermissions(ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	var org *models.Organization
	if in.GetOrganizationId() != 0 {
		org = &models.Organization{ID: int(in.GetOrganizationId())}
	}

	// TODO валидация на сложный пароль
	if in.Password != "" && len(in.Password) < 5 {
		log.Println("Failed to validate password")
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	// Если указан пароль, хешируем его перед сохранением
	var passwordHash string
	if in.Password != "" {
		hash, err := utils.HashPassword(in.Password)
		if err != nil {
			log.Printf("Failed to hash password, err:%v\n", err)
			return nil, status.Error(codes.Internal, "Failed to update user")
		}
		passwordHash = hash
	}

	// Преобразуем запрос в модель
	user := models.User{
		ID:           int(in.Id),
		Name:         in.Name,
		Email:        in.Email,
		PasswordHash: passwordHash,
		Organization: org,
	}

	// Используем репозиторий для обновления пользователя
	if err := bh.userRepo.UpdateUser(ctx, &user); err != nil {
		log.Printf("update user failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "Failed to update user")
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

	// Возвращаем ответ
	return &api_pb.User{
		Id:           int32(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		Organization: orgOut,
	}, nil
}

func checkPermissions(ctx context.Context, id int32) error {
	userID, ok := ctx.Value(auth.UserIDKey).(int)
	if !ok {
		return errors.New("auth.UserIDKey")
	}

	isAdmin, ok := ctx.Value(auth.IsAdminKey).(bool)
	if !ok {
		return errors.New("auth.IsAdminKey")
	}

	if userID == int(id) || isAdmin {
		return nil
	}

	return errors.New("checkPermissions.Unauthenticated")
}
