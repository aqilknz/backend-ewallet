package repository

import (
	"context"

	"github.com/aqilknz/backend-ewallet/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, tx pgx.Tx, email, password string) (model.User, error)
	CreateProfile(ctx context.Context, tx pgx.Tx, userID uint) error
	CreateWallet(ctx context.Context, tx pgx.Tx, userID uint) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, sql, email).Scan(&exists)
	return exists, err
}

func (r *authRepository) CreateUser(ctx context.Context, tx pgx.Tx, email, password string) (model.User, error) {
	sql := `INSERT INTO users (email, password, pin, created_at, updated_at) VALUES ($1, $2, '', NOW(), NOW()) RETURNING id, email, created_at, updated_at`
	args := []any{email, password}

	var user model.User
	err := tx.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (r *authRepository) CreateProfile(ctx context.Context, tx pgx.Tx, userID uint) error {
	sql := `INSERT INTO profiles (user_id, full_name, phone, photo) VALUES ($1, '', '', '')`
	args := []any{userID}
	_, err := tx.Exec(ctx, sql, args...)
	return err
}

func (r *authRepository) CreateWallet(ctx context.Context, tx pgx.Tx, userID uint) error {
	sql := `INSERT INTO wallets (user_id, balance) VALUES ($1, 0)`
	args := []any{userID}
	_, err := tx.Exec(ctx, sql, args...)
	return err
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	sql := `SELECT id, email, password FROM users WHERE email = $1`
	args := []any{email}

	var user model.User
	err := r.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Email, &user.Password)
	return user, err
}
