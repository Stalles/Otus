package getuserbyid

import (
	"context"
	"errors"
	"socialNetworkOtus/internal/api"
	"socialNetworkOtus/internal/repository"
)

type Service struct {
	UserRepo *repository.UserRepository
}

func NewService(userRepo *repository.UserRepository) *Service {
	return &Service{UserRepo: userRepo}
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*api.User, error) {
	user, err := s.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found: " + err.Error())
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
