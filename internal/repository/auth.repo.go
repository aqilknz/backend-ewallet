package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aqilknz/backend-ewallet/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthRepository interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, tx pgx.Tx, email, password string) (model.User, error)
	CreateProfile(ctx context.Context, tx pgx.Tx, userID uint) error
	CreateWallet(ctx context.Context, tx pgx.Tx, userID uint) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	AddTokenToBlacklist(ctx context.Context, userID int, token string, expiresIn time.Duration) error
	IsTokenBlacklisted(ctx context.Context, userID int, token string) bool
	CreatePin(ctx context.Context, userID int, pinHash string) error
	GetUserPin(ctx context.Context, userID int) (string, error)
	UpdatePassword(ctx context.Context, email string, hashedPassword string) error
}

type authRepository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewAuthRepository(db *pgxpool.Pool, redis *redis.Client) AuthRepository {
	return &authRepository{
		db:    db,
		redis: redis,
	}
}

func (ar *authRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := ar.db.QueryRow(ctx, sql, email).Scan(&exists)
	return exists, err
}

func (ar *authRepository) CreateUser(ctx context.Context, tx pgx.Tx, email, password string) (model.User, error) {
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

func (ar *authRepository) CreateProfile(ctx context.Context, tx pgx.Tx, userID uint) error {
	sql := `INSERT INTO profiles (user_id, full_name, phone, photo) VALUES ($1, '', '', '')`
	args := []any{userID}
	_, err := tx.Exec(ctx, sql, args...)
	return err
}

func (ar *authRepository) CreateWallet(ctx context.Context, tx pgx.Tx, userID uint) error {
	sql := `INSERT INTO wallets (user_id, balance) VALUES ($1, 0)`
	args := []any{userID}
	_, err := tx.Exec(ctx, sql, args...)
	return err
}

func (ar *authRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	sql := `SELECT id, email, password, pin FROM users WHERE email = $1`
	args := []any{email}

	var user model.User
	err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Email, &user.Password, &user.Pin)
	if err == pgx.ErrNoRows {
		return user, errors.New("user not found")
	}
	return user, err
}

func (ar *authRepository) AddTokenToBlacklist(ctx context.Context, userID int, token string, expiresIn time.Duration) error {
	key := fmt.Sprintf("blacklist:user:%d:token:%s", userID, token)
	return ar.redis.Set(ctx, key, "revoked", expiresIn).Err()
}

func (r *authRepository) IsTokenBlacklisted(ctx context.Context, userID int, token string) bool {
	key := fmt.Sprintf("blacklist:user:%d:token:%s", userID, token)
	err := r.redis.Get(ctx, key).Err()

	if err == redis.Nil {
		return false
	}

	return true
}

func (ar *authRepository) GetUserPin(ctx context.Context, userID int) (string, error) {
	sql := `SELECT pin FROM users WHERE id = $1`
	var pin string
	err := ar.db.QueryRow(ctx, sql, userID).Scan(&pin)
	return pin, err
}

func (ar *authRepository) CreatePin(ctx context.Context, userID int, pinHash string) error {
	sql := `UPDATE users SET pin = $1, updated_at = NOW() WHERE id = $2`
	_, err := ar.db.Exec(ctx, sql, pinHash, userID)
	return err
}

func (ar *authRepository) UpdatePassword(ctx context.Context, email string, hashedPassword string) error {
	sqlQuery := `UPDATE users SET password = $1, updated_at = NOW() WHERE email = $2`

	result, err := ar.db.Exec(ctx, sqlQuery, hashedPassword, email)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("gagal menyimpan perubahan ke database (email tidak ditemukan)")
	}

	return nil
}
