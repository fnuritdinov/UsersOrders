package service

import (
	"UserService/internal/models"
	"UserService/internal/repository"
	errs "UserService/pkg/errors"
	"context"
	"fmt"
)

type AdminService interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	ChangeRole(ctx context.Context, id int, user models.User) (models.User, error)
}

type serviceAdmin struct {
	RepoAdmin repository.AdminRepo
}

func NewAdminService(RepoAdmin repository.AdminRepo) AdminService {
	return &serviceAdmin{
		RepoAdmin: RepoAdmin,
	}
}

func (s *serviceAdmin) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.RepoAdmin.GetAllUsers(ctx)
	if err != nil {
		return []models.User{}, fmt.Errorf("error from s.RepoAdmin.GetAllUsers")
	}
	return users, nil
}

func (s *serviceAdmin) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	orders, err := s.RepoAdmin.GetAllOrders(ctx)
	if err != nil {
		return []models.Order{}, nil
	}

	return orders, nil
}

func (s *serviceAdmin) ChangeRole(ctx context.Context, id int, user models.User) (models.User, error) {
	if user.Role != models.AdminRole && user.Role != models.UserRole {
		return models.User{}, errs.ErrValidate
	}

	user, err := s.RepoAdmin.ChangeRole(ctx, id, user)
	if err != nil {
		return models.User{}, fmt.Errorf("error from s.RepoAdmin.ChangeRole")
	}
	return user, nil
}
