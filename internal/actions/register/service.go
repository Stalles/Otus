package register

import (
	"context"
	"fmt"
	"socialNetworkOtus/internal/api"
	"socialNetworkOtus/internal/repository"
	"socialNetworkOtus/internal/utils"
)

type Service struct {
	UserRepo *repository.UserRepository
}

func NewService(userRepo *repository.UserRepository) *Service {
	return &Service{UserRepo: userRepo}
}

func (s *Service) RegisterUser(ctx context.Context, req *api.PostUserRegisterJSONBody) (string, error) {
	passwordHash, err := utils.HashPassword(*req.Password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	userID, err := s.UserRepo.CreateUser(ctx, req, passwordHash)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	return userID, nil
}
