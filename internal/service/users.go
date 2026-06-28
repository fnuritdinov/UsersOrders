package service

import (
	"UserService/internal/models"
	"UserService/internal/repository"
	errs "UserService/pkg/errors"
	"UserService/pkg/jwt"
	"UserService/pkg/memory"
	"UserService/pkg/password"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type UserService interface {
	Register(ctx context.Context, request models.RegisterRequest) error
	Login(ctx context.Context, login models.LoginRequest) (string, error)
	Verify(ctx context.Context, request models.RegisterRequest) (int, error)

	Get(ctx context.Context, userID int) (models.User, error)
	UpdateProfile(ctx context.Context, userID int, user models.User) error
	DeleteProfile(ctx context.Context, userID int) error
	ChangePassword(ctx context.Context, userID int, user models.Password) error
	GetOrdersProfile(ctx context.Context, userID int) (models.UserOrder, error)
}

type serviceUser struct {
	RepoUser  repository.UserRepo
	RepoOrder repository.OrderRepo
	myCache   memory.MemoryCache
}

func NewUserService(RepoUser repository.UserRepo, RepoOrder repository.OrderRepo, myCache memory.MemoryCache) UserService {
	return &serviceUser{
		RepoUser:  RepoUser,
		RepoOrder: RepoOrder,
		myCache:   myCache,
	}
}

func (s *serviceUser) Register(ctx context.Context, request models.RegisterRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	exists, err := s.RepoUser.ExistsByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("user with this email already exists")
	}

	passwordHash, err := password.Hash(request.Password)
	if err != nil {
		return err
	}

	s.myCache.Set(request.Email, models.RegisterRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: passwordHash,
		Role:     models.UserRole,
	}, 5*time.Minute)

	return nil
}

type CacheMemory struct {
	Name     string
	Email    string
	Password string
	Role     string
	OTP      string
}

func (s *serviceUser) Verify(ctx context.Context, request models.RegisterRequest) (int, error) {
	data, ok := s.myCache.Get(request.Email)
	if !ok {
		return 0, errors.New("user not found in cache")
	}

	cacheInfo, ok := data.(CacheMemory)
	if !ok {
		return 0, errors.New("internal error")
	}

	id, err := s.RepoUser.Register(ctx, models.RegisterRequest{
		Name:     cacheInfo.Name,
		Email:    cacheInfo.Email,
		Password: cacheInfo.Password,
		Role:     cacheInfo.Role,
	})
	if err != nil {
		return 0, err
	}

	s.myCache.Remove(request.Email)

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

	err = password.Compare(user.Password, request.Password)
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

func (s *serviceUser) ChangePassword(ctx context.Context, userID int, request models.Password) error {
	userInfo, err := s.RepoUser.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("error from s.RepoUser.Get %w", err)
	}

	err = password.Compare(userInfo.Password, request.OldPassword)
	if err != nil {
		return fmt.Errorf("error from password.Campare")
	}

	hashPassword, err := password.Hash(request.NewPassword)
	if err != nil {
		return fmt.Errorf("error from password.Hash %w", err)
	}

	err = s.RepoUser.ChangePassword(ctx, userID, hashPassword)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrNotFound
		}
		return err
	}

	return nil
}

func (s *serviceUser) GetOrdersProfile(ctx context.Context, userID int) (models.UserOrder, error) {
	user, err := s.RepoUser.Get(ctx, userID)
	if err != nil {
		return models.UserOrder{}, err
	}

	orders, err := s.RepoOrder.GetOrder(ctx, userID)
	if err != nil {
		return models.UserOrder{}, err
	}
	return models.UserOrder{
		User:   user,
		Orders: orders,
	}, nil
}
