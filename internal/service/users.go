package service

import (
	"UserService/internal/models"
	"UserService/internal/repository"
	errs "UserService/pkg/errors"
	"UserService/pkg/jwt"
	"UserService/pkg/password"
	"context"
	"fmt"
	"log"
)

type UserService interface {
	Register(ctx context.Context, request models.RegisterRequest) (int, error)
	Login(ctx context.Context, login models.LoginRequest) (string, error)

	Get(ctx context.Context, userID int) (models.User, error)
	UpdateProfile(ctx context.Context, userID int, user models.User) error
	DeleteProfile(ctx context.Context, userID int) error
}

type serviceUser struct {
	RepoUser repository.UserRepo
}

func NewUserService(RepoUser repository.UserRepo) UserService {
	return &serviceUser{
		RepoUser: RepoUser,
	}
}

func (s *serviceUser) Register(ctx context.Context, request models.RegisterRequest) (int, error) {
	err := request.Validate()
	if err != nil {
		return 0, err
	}

	exists, err := s.RepoUser.ExistsByEmail(ctx, request.Email)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, fmt.Errorf("user with this email already exists")
	}

	passwordHash, err := password.Hash(request.Password)
	if err != nil {
		return 0, err
	}

	id, err := s.RepoUser.Register(ctx, models.RegisterRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: passwordHash,
		Role:     models.UserRole,
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *serviceUser) Login(ctx context.Context, request models.LoginRequest) (string, error) {
	err := request.Validate()
	if err != nil {
		return "", err
	}

	user, err := s.RepoUser.GetByEmail(ctx, request.Email)
	if err != nil {
		return "", err
	}

	err = password.Campare(user.Password, request.Password)
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("error from jwt.GenerateToken %w", err)
	}

	return token, err
}

func (s *serviceUser) Get(ctx context.Context, userID int) (models.User, error) {

	user, err := s.RepoUser.Get(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *serviceUser) UpdateProfile(ctx context.Context, userID int, user models.User) error {

	err := models.StrEmpty(user.Name)
	if err != nil {
		return err
	}

	log.Println(user.Name)

	if !models.IsValidPhone(user.Phone) {
		return errs.ErrValidate
	}

	log.Println(user.Phone)

	err = s.RepoUser.UpdateProfile(ctx, userID, user)
	if err != nil {
		return fmt.Errorf("error from s.Repo.UpdateProfile %w", err)
	}

	return nil
}

func (s *serviceUser) DeleteProfile(ctx context.Context, userID int) error {

	err := s.RepoUser.DeleteProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("error from s.Repo.DeleteProfile %w", err)
	}
	return nil
}
