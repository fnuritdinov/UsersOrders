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
	CancelOrder(ctx context.Context, userID int, orderID int) error
}

type serviceOrder struct {
	repoOrder repository.OrderRepo
}

func NewOrderService(repoOrder repository.OrderRepo) OrderService {
	return &serviceOrder{
		repoOrder: repoOrder,
	}
}

func (s *serviceOrder) CreateOrder(ctx context.Context, order models.Order) error {
	err := order.ValidateCreateOrder()
	if err != nil {
		return err
	}

	err = s.repoOrder.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceOrder) GetMyOrders(ctx context.Context, userID int) ([]models.Order, error) {

	orders, err := s.repoOrder.GetMyOrders(ctx, userID)
	if err != nil {
		return []models.Order{}, fmt.Errorf("error from s.Repo.GetMyOrders %w", err)
	}

	return orders, nil
}

func (s *serviceOrder) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	order, err := s.repoOrder.GetOrderByID(ctx, id)
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

	if updatedOrder.Price < 1 {
		return errs.ErrValidate
	}

	err = s.repoOrder.UpdateOrder(ctx, id, updatedOrder)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceOrder) DeleteOrders(ctx context.Context, id int) error {

	err := s.repoOrder.DeleteOrders(ctx, id)
	if err != nil {
		return fmt.Errorf("error from s.Repo.DeleteOrders %w", err)
	}

	return nil
}

func (s *serviceOrder) CancelOrder(ctx context.Context, userID int, orderID int) error {

	order, err := s.repoOrder.GetOrder(ctx, userID)
	if err != nil {
		return err
	}

	if order.Status != "new" {
		return fmt.Errorf("invalid status for cancel")
	}

	err = s.repoOrder.CancelOrder(ctx, orderID)
	if err != nil {
		return err
	}
	return nil
}
