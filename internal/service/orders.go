package service

import (
	"UserService/internal/models"
	"UserService/internal/repository"
	errs "UserService/pkg/errors"
	"context"
	"fmt"
	"log"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetMyOrders(ctx context.Context, userID int) ([]models.Order, error)
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
	UpdateOrder(ctx context.Context, id int, updatedOrder models.Order) error
	DeleteOrders(ctx context.Context, id int) error
}

type serviceOrder struct {
	RepoOrder repository.OrderRepo
}

func NewOrderService(RepoOrder repository.OrderRepo) OrderService {
	return &serviceOrder{
		RepoOrder: RepoOrder,
	}
}

func (s *serviceOrder) CreateOrder(ctx context.Context, order models.Order) error {
	err := order.ValidateCreateOrder()
	if err != nil {
		return err
	}

	err = s.RepoOrder.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceOrder) GetMyOrders(ctx context.Context, userID int) ([]models.Order, error) {

	orders, err := s.RepoOrder.GetMyOrders(ctx, userID)
	if err != nil {
		return []models.Order{}, fmt.Errorf("error from s.Repo.GetMyOrders %w", err)
	}

	return orders, nil
}

func (s *serviceOrder) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	order, err := s.RepoOrder.GetOrderByID(ctx, id)
	if err != nil {
		return models.Order{}, fmt.Errorf("error from s.Repo.GetOrderByID %w", err)
	}
	return order, nil
}

func (s *serviceOrder) UpdateOrder(ctx context.Context, id int, updatedOrder models.Order) error {
	err := models.StrEmpty(updatedOrder.Product)
	if err != nil {
		return err
	}

	log.Println(updatedOrder.Product)

	err = models.StrEmpty(updatedOrder.Status)
	if err != nil {
		return err
	}
	log.Println(updatedOrder.Status)

	if updatedOrder.Price < 1 {
		return errs.ErrValidate
	}
	log.Println(updatedOrder.Price)

	err = s.RepoOrder.UpdateOrder(ctx, id, updatedOrder)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceOrder) DeleteOrders(ctx context.Context, id int) error {

	err := s.RepoOrder.DeleteOrders(ctx, id)
	if err != nil {
		return fmt.Errorf("error from s.Repo.DeleteOrders")
	}

	return nil
}
