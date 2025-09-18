package service

import (
	"context"
	"errors"
	"event-booking-system/internal/config"
	"event-booking-system/internal/domain"
	"event-booking-system/internal/models"
	"event-booking-system/internal/repository"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var (
	jwtSecret = config.GetConfig().JWTSecret
	jwtExp    = 3 * time.Hour
)

type UserService struct {
	repo  *repository.UserRepository
	cache *redis.Client
}

func NewUserService(repo *repository.UserRepository, cache *redis.Client) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

func (s *UserService) Create(ctx context.Context, payload *models.RequestCreateUser) error {
	var user models.User
	user.Name = payload.Name
	user.Email = payload.Email
	user.Password = payload.Password

	err := s.repo.Create(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, user *models.RequestUpdateUser) error {
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Save user ID to redis denylist
	s.cache.Set(ctx, fmt.Sprintf("denylist:%s", id), id, jwtExp)

	return nil
}

func (s *UserService) Register(ctx context.Context, payload models.RequestRegisterUser) error {
	var user models.User
	user.Name = payload.Name
	user.Email = payload.Email
	user.Password = payload.Password

	err := s.repo.Create(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, payload models.RequestLoginUser) (string, error) {
	email := payload.Email
	password := payload.Password

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("User not found")
	}

	if user.Password != password {
		return "", errors.New("Credential invalid")
	}

	claims := domain.UserClaims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
