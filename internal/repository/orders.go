package repository

import (
	"UserService/internal/models"
	errs "UserService/pkg/errors"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetMyOrders(ctx context.Context, userID int) ([]models.Order, error)
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
	UpdateOrder(ctx context.Context, id int, updateOrder models.Order) error
	DeleteOrders(ctx context.Context, id int) error
}

type repoOrder struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) OrderRepo {
	return &repoOrder{
		db: db,
	}
}

func (r *repoOrder) CreateOrder(ctx context.Context, order models.Order) error {
	const query = `
		INSERT INTO orders (product, price, user_id, status)
		VALUES ($1, $2, $3, $4)`

	result, err := r.db.Exec(ctx, query, order.Product, order.Price, order.UserID, order.Status)
	if err != nil {
		return fmt.Errorf("error from r.Pgx.Exec %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order was not created")
	}

	return nil
}

func (r *repoOrder) GetMyOrders(ctx context.Context, userID int) ([]models.Order, error) {
	const query = `SELECT id, product, price, status FROM orders WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return []models.Order{}, fmt.Errorf("error from r.Pgx.Query %w", err)
	}

	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order

		err = rows.Scan(
			&order.ID,
			&order.Product,
			&order.Price,
			&order.Status)
		if err != nil {
			return nil, fmt.Errorf("error from rows.Scan %w", err)
		}

		orders = append(orders, order)

	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error from rows.Err %w", err)
	}

	return orders, nil
}

func (r *repoOrder) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	var order models.Order

	const query = `
			SELECT id, product, price, user_id, status, created_at
			FROM orders
			WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.Product,
		&order.Price,
		&order.UserID,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Order{}, errs.ErrNotFound
		}
		return models.Order{}, fmt.Errorf("error from r.Pgx.QueryRow %w", err)
	}
	return order, nil
}

func (r *repoOrder) UpdateOrder(ctx context.Context, id int, updateOrder models.Order) error {
	const query = `
		UPDATE orders 
		SET  product = $2, price = $3, status = $4 
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, updateOrder.Product, updateOrder.Price, updateOrder.Status)
	if err != nil {
		return fmt.Errorf("error from r.Pgx.Exec %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *repoOrder) DeleteOrders(ctx context.Context, id int) error {

	const query = `DELETE FROM orders WHERE id = $1`

	rows, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error from delete %w", err)
	}

	if rows.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
