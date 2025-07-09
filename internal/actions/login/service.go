package login

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"socialNetworkOtus/internal/repository"
	"socialNetworkOtus/internal/utils"

	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v4"
)

type Service struct {
	UserRepo  *repository.UserRepository
	JWTSecret string
}

func NewService(userRepo *repository.UserRepository) *Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret"
	}
	return &Service{UserRepo: userRepo, JWTSecret: secret}
}

func (s *Service) LoginUser(ctx context.Context, id, password string) (string, error) {
	user, err := s.UserRepo.GetUserByID(ctx, id)
	if err != nil || user == nil {
		return "", errors.New("user not found")
	}

	passwordHash, err := s.getPasswordHash(ctx, id)
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(password, passwordHash) {
		return "", errors.New("invalid password")
	}

	return s.generateToken(id)
}

func (s *Service) getPasswordHash(ctx context.Context, id string) (string, error) {
	var passwordHash string
	db := s.UserRepo.DB()
	ds := db.From("users").Select("password_hash").Where(goqu.Ex{"id": id})
	found, err := ds.ScanValContext(ctx, &passwordHash)
	if err != nil {
		return "", fmt.Errorf("failed to get password hash: %w", err)
	}
	if !found {
		return "", errors.New("user not found")
	}
	return passwordHash, nil
}

func (s *Service) generateToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tokenString, nil
}
