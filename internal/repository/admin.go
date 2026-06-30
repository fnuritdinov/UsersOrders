package repository

import (
	"UserService/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepo interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	ChangeRole(ctx context.Context, id int, user models.User) (models.User, error)
}

type repoAdmin struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) AdminRepo {
	return &repoAdmin{
		db: db,
	}
}

func (r *repoAdmin) GetAllUsers(ctx context.Context) ([]models.User, error) {

	const query = `
			SELECT id, name, email, phone, password_hash, role 
			FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []models.User{}, fmt.Errorf("error from r.db.Query")
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Password,
			&user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repoAdmin) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	const query = `
		SELECT id, product, price, user_id, status 
		FROM orders`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []models.Order{}, fmt.Errorf("error from r.db.Query")
	}

	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order

		err = rows.Scan(
			&order.ID,
			&order.Product,
			&order.Price,
			&order.UserID,
			&order.Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *repoAdmin) ChangeRole(ctx context.Context, id int, request models.User) (models.User, error) {
	const query = `UPDATE users SET role = $2 WHERE id = $1 RETURNING id, name, email, phone, role`

	var user models.User

	err := r.db.QueryRow(ctx, query, id, request.Role).
		Scan(&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Role)

	if err != nil {
		return models.User{}, fmt.Errorf("error from r.db.QueryRow %w", err)
	}

	return user, nil
}
