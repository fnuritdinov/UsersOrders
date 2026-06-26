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

type UserRepo interface {
	Register(ctx context.Context, request models.RegisterRequest) (int, error)
	Login(ctx context.Context, request models.LoginRequest) error

	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)

	Get(ctx context.Context, userID int) (models.User, error)
	UpdateProfile(ctx context.Context, userID int, user models.User) error
	DeleteProfile(ctx context.Context, userID int) error
}

type repoUser struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepo {
	return &repoUser{
		db: db,
	}
}

func (r *repoUser) Register(ctx context.Context, request models.RegisterRequest) (int, error) {

	const query = `
			INSERT INTO users (name, email, password_hash, role)
			VALUES ($1, $2, $3, $4) RETURNING id;`

	var id int

	err := r.db.QueryRow(ctx, query, request.Name, request.Email, request.Password, request.Role).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error from r.Pgx.QueryRow: %w", err)
	}

	return id, nil
}

func (r *repoUser) Login(ctx context.Context, request models.LoginRequest) error {
	return nil
}

func (r *repoUser) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `
			SELECT EXISTS (
				SELECT 1 
				FROM users
				WHERE email = $1
			);`

	var exist bool

	err := r.db.QueryRow(ctx, query, email).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (r *repoUser) GetByEmail(ctx context.Context, email string) (models.User, error) {
	const query = `
			SELECT id, name, email, password_hash, role 
			FROM users
			WHERE email = $1`

	var user models.User

	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *repoUser) Get(ctx context.Context, userID int) (models.User, error) {
	var user models.User
	const query = `
			SELECT id, name, email, password_hash, role 
			FROM users
			WHERE id = $1`

	err := r.db.QueryRow(ctx, query, userID).
		Scan(&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errs.ErrNotFound
		}
		return models.User{}, fmt.Errorf("error from r.Pgx.QueryRow %w", err)
	}
	return user, nil
}

func (r *repoUser) UpdateProfile(ctx context.Context, userID int, user models.User) error {
	const query = `UPDATE users SET name = $2, phone = $3 WHERE id = $1`

	row, err := r.db.Exec(ctx, query, userID, user.Name, user.Phone)
	if err != nil {
		return fmt.Errorf("error from update %w", err)
	}

	if row.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *repoUser) DeleteProfile(ctx context.Context, userID int) error {

	const query = `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error from delete user %w", err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
