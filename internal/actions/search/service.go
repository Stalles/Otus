package search

import (
	"context"
	"socialNetworkOtus/internal/api"
	"socialNetworkOtus/internal/repository"
)

type Service struct {
	UserRepo *repository.UserRepository
}

func NewService(userRepo *repository.UserRepository) *Service {
	return &Service{UserRepo: userRepo}
}

func (s *Service) SearchUsersByPrefix(ctx context.Context, firstName, lastName string) ([]api.User, error) {
	return s.UserRepo.SearchUsersByPrefix(ctx, firstName, lastName)
}
