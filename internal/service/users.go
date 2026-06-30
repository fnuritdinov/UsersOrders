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
	"time"
)

type UserService interface {
	Register(ctx context.Context, request models.RegisterRequest) error
	Login(ctx context.Context, login models.LoginRequest) (string, string, error)
	Verify(ctx context.Context, request models.RegisterRequest) (int, error)
	RefreshToken(ctx context.Context, require models.HashToken) (models.RefreshAccessTokens, error)
	LogOut(ctx context.Context, request models.RefreshAccessTokens) error

	Get(ctx context.Context, userID int) (models.User, error)
	UpdateProfile(ctx context.Context, userID int, user models.User) error
	DeleteProfile(ctx context.Context, userID int) error
	ChangePassword(ctx context.Context, userID int, user models.Password) error
	GetOrdersProfile(ctx context.Context, userID int) (models.UserOrder, error)
}

type serviceUser struct {
	repoUser  repository.UserRepo
	repoOrder repository.OrderRepo
	myCache   memory.MemoryCache
}

func NewUserService(repoUser repository.UserRepo, repoOrder repository.OrderRepo, myCache memory.MemoryCache) UserService {
	return &serviceUser{
		repoUser:  repoUser,
		repoOrder: repoOrder,
		myCache:   myCache,
	}
}

func (s *serviceUser) Register(ctx context.Context, request models.RegisterRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	exists, err := s.repoUser.ExistsByEmail(ctx, request.Email)
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

	s.myCache.Set(request.Email, CacheMemory{
		Name:     request.Name,
		Email:    request.Email,
		Password: passwordHash,
		OTP:      "12345",
		Role:     models.UserRole,
	}, 5*time.Minute)

	return nil
}

type CacheMemory struct {
	Name        string
	Email       string
	Password    string
	Role        string
	OTP         string
	AttemptInfo int
}

func (s *serviceUser) Verify(ctx context.Context, request models.RegisterRequest) (int, error) {
	data, ok := s.myCache.Get(request.Email)
	if !ok {
		return 0, errors.New("user not found in cache")
	}

	cacheInfo, ok := data.(CacheMemory)
	if !ok {
		return 0, errors.New("data.(CacheMemory)")
	}

	if cacheInfo.AttemptInfo >= 3 {
		return 0, errors.New("user is too many attempts")
	}

	if request.OTP != cacheInfo.OTP {
		cacheInfo.AttemptInfo++
		s.myCache.Set(request.Email, cacheInfo, time.Minute*5)
		return 0, errors.New("invalid otp")
	}

	id, err := s.repoUser.Register(ctx, models.RegisterRequest{
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

func (s *serviceUser) Login(ctx context.Context, request models.LoginRequest) (string, string, error) {
	err := request.Validate()
	if err != nil {
		return "", "", err
	}

	user, err := s.repoUser.GetByEmail(ctx, request.Email)
	if err != nil {
		return "", "", err
	}

	err = password.Compare(user.Password, request.Password)
	if err != nil {
		return "", "", err
	}

	accessToken, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", "", fmt.Errorf("error from jwt.GenerateToken %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", errors.New("error from jwt.GenerateRefreshToken")
	}

	hash := jwt.HashRefreshToken(refreshToken)

	err = s.repoUser.SaveRefreshToken(ctx, models.HashToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(168 * time.Hour),
	})

	return accessToken, refreshToken, nil
}

func (s *serviceUser) Get(ctx context.Context, userID int) (models.User, error) {

	user, err := s.repoUser.Get(ctx, userID)
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

	if !models.IsValidPhone(user.Phone) {
		return errs.ErrValidate
	}

	err = s.repoUser.UpdateProfile(ctx, userID, user)
	if err != nil {
		return fmt.Errorf("error from s.Repo.UpdateProfile %w", err)
	}

	return nil
}

func (s *serviceUser) DeleteProfile(ctx context.Context, userID int) error {

	err := s.repoUser.DeleteProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("error from s.Repo.DeleteProfile %w", err)
	}
	return nil
}

func (s *serviceUser) ChangePassword(ctx context.Context, userID int, request models.Password) error {
	userInfo, err := s.repoUser.Get(ctx, userID)
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

	err = s.repoUser.ChangePassword(ctx, userID, hashPassword)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrNotFound
		}
		return err
	}

	return nil
}

func (s *serviceUser) GetOrdersProfile(ctx context.Context, userID int) (models.UserOrder, error) {
	user, err := s.repoUser.Get(ctx, userID)
	if err != nil {
		return models.UserOrder{}, err
	}

	orders, err := s.repoOrder.GetOrder(ctx, userID)
	if err != nil {
		return models.UserOrder{}, err
	}
	return models.UserOrder{
		User:   user,
		Orders: orders,
	}, nil
}

func (s *serviceUser) RefreshToken(ctx context.Context, request models.HashToken) (models.RefreshAccessTokens, error) {
	hash := jwt.HashRefreshToken(request.TokenHash)

	token, err := s.repoUser.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return models.RefreshAccessTokens{}, err
	}

	user, err := s.repoUser.Get(ctx, token.UserID)
	if err != nil {
		return models.RefreshAccessTokens{}, err
	}

	if time.Now().After(token.ExpiresAt) {
		return models.RefreshAccessTokens{}, errs.ErrUnauthorized
	}

	err = s.repoUser.DeleteRefreshTokenByID(ctx, token.ID)
	if err != nil {
		return models.RefreshAccessTokens{}, err
	}

	newRefreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return models.RefreshAccessTokens{}, err
	}

	newHash := jwt.HashRefreshToken(newRefreshToken)

	err = s.repoUser.SaveRefreshToken(ctx, models.HashToken{
		UserID:    token.UserID,
		TokenHash: newHash,
		ExpiresAt: time.Now().Add(time.Hour * 168),
	})
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return models.RefreshAccessTokens{}, errs.ErrNotFound
		}
		return models.RefreshAccessTokens{}, err
	}

	accessToken, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.RefreshAccessTokens{}, err
	}

	return models.RefreshAccessTokens{
		RefreshToken: newRefreshToken,
		AccessToken:  accessToken,
	}, nil
}

func (s *serviceUser) LogOut(ctx context.Context, request models.RefreshAccessTokens) error {
	hash := jwt.HashRefreshToken(request.RefreshToken)

	err := s.repoUser.DeleteRefreshToken(ctx, hash)
	if err != nil {
		return err
	}
	return nil
}
