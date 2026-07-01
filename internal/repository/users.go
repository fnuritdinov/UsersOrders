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
	LoginHistory(ctx context.Context, userID int, ip, userAgent string) error
	GetLoginHistory(ctx context.Context, userID int) ([]models.LoginHistoryResponse, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)

	Get(ctx context.Context, userID int) (models.User, error)
	UpdateProfile(ctx context.Context, userID int, user models.User) error
	DeleteProfile(ctx context.Context, userID int) error
	ChangePassword(ctx context.Context, userID int, hashPassword string) error
	SaveRefreshToken(ctx context.Context, request models.HashToken) error
	GetRefreshTokenByHash(ctx context.Context, hash string) (models.HashToken, error)
	DeleteRefreshTokenByID(ctx context.Context, tokenID int) error
	DeleteRefreshToken(ctx context.Context, token string) error
	GetRefreshTokenByUserID(ctx context.Context, userID int) (models.HashToken, error)
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

func (r *repoUser) LoginHistory(ctx context.Context, userID int, ip, userAgent string) error {
	const query = `INSERT INTO login_history (user_id, ip, user_agent) VALUES ($1, $2, $3)`

	result, err := r.db.Exec(ctx, query, userID, ip, userAgent)
	if err != nil {
		return errors.New("error from r.db.Exec")
	}

	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *repoUser) GetLoginHistory(ctx context.Context, userID int) ([]models.LoginHistoryResponse, error) {
	const query = `
			SELECT ip, user_agent, created_at 
			FROM login_history
			WHERE user_id = $1
			ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var histories []models.LoginHistoryResponse

	for rows.Next() {
		var history models.LoginHistoryResponse

		if err = rows.Scan(
			&history.IP,
			&history.UserAgent,
			&history.CreatedAt); err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return histories, nil
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

func (r *repoUser) ChangePassword(ctx context.Context, userID int, hashPassowrd string) error {
	const query = `UPDATE users SET password_hash = $2 WHERE id = $1`

	rows, err := r.db.Exec(ctx, query, userID, hashPassowrd)
	if err != nil {
		return fmt.Errorf("error from r.db.Exec %w", err)
	}

	if rows.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *repoUser) SaveRefreshToken(ctx context.Context, request models.HashToken) error {
	const query = `INSERT INTO refresh_tokens(user_id, token_hash, expires_at) VALUES ($1, $2, $3)`

	rows, err := r.db.Exec(ctx, query, request.UserID, request.TokenHash, request.ExpiresAt)
	if err != nil {
		return errors.New("error from r.db.Exec")
	}

	if rows.RowsAffected() == 0 {
		return errors.New("token was not created")
	}
	return nil
}

func (r *repoUser) GetRefreshTokenByHash(ctx context.Context, hash string) (models.HashToken, error) {
	var token models.HashToken
	const query = `SELECT id, user_id, expires_at FROM refresh_tokens WHERE token_hash = $1`

	err := r.db.QueryRow(ctx, query, hash).
		Scan(&token.ID,
			&token.UserID,
			&token.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.HashToken{}, errs.ErrNotFound
		}

		return models.HashToken{}, fmt.Errorf("error from r.db.QueryRow %w", err)
	}

	return token, nil
}

func (r *repoUser) DeleteRefreshTokenByID(ctx context.Context, tokenID int) error {
	const query = `DELETE FROM refresh_tokens WHERE id = $1`

	result, err := r.db.Exec(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("error from delete token %w", err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (r *repoUser) DeleteRefreshToken(ctx context.Context, token string) error {
	const query = `DELETE FROM refresh_tokens WHERE token_hash = $1`

	result, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("error from r.db.Exec", err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *repoUser) GetRefreshTokenByUserID(ctx context.Context, userID int) (models.HashToken, error) {
	const query = `SELECT token_hash FROM refresh_tokens WHERE user_id = $1`

	var token models.HashToken
	err := r.db.QueryRow(ctx, query, userID).Scan(&token.TokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.HashToken{}, errs.ErrNotFound
		}
		return models.HashToken{}, errors.New("error from r.db.QueryRow")
	}
	return token, nil
}
