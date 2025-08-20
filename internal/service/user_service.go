package service

import (
	"context"
	"errors"
	"event-booking-system/internal/config"
	"event-booking-system/internal/domain"
	"event-booking-system/internal/models"
	"event-booking-system/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = config.GetConfig().JWTSecret
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) Create(ctx context.Context, payload *models.RequestCreateUser) error {
	var user models.User
	user.Name = payload.Name
	user.Email = payload.Email
	user.Password = payload.Password

	err := us.repo.Create(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) List(ctx context.Context) ([]models.User, error) {
	users, err := us.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, err := us.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Update(ctx context.Context, user *models.RequestUpdateUser) error {
	if err := us.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (us *UserService) Delete(ctx context.Context, id string) error {
	if err := us.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (us *UserService) Register(ctx context.Context, payload models.RequestRegisterUser) error {
	var user models.User
	user.Name = payload.Name
	user.Email = payload.Email
	user.Password = payload.Password

	err := us.repo.Create(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Login(ctx context.Context, payload models.RequestLoginUser) (string, error) {
	email := payload.Email
	password := payload.Password

	user, err := us.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("User not found")
	}

	if user.Password != password {
		return "", errors.New("Credential invalid")
	}

	claims := domain.UserClaims{
		ID:   user.ID,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
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
