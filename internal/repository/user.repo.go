package repository

import (
	"context"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/aqilknz/backend-ewallet/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetProfile(ctx context.Context, userID int) (dto.UserProfileResponse, error)
	GetDashboard(ctx context.Context, userID int) (dto.DashboardResponse, error)
	UpdateProfile(ctx context.Context, userID int, req dto.EditProfileRequest) error
	GetPasswordAndPin(ctx context.Context, userID int) (password string, pin string, err error)
	UpdatePassword(ctx context.Context, userID int, newPassword string) error
	UpdatePin(ctx context.Context, userID int, newPin string) error
	GetUserByID(ctx context.Context, userID int) (model.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetProfile(ctx context.Context, userID int) (dto.UserProfileResponse, error) {
	var res dto.UserProfileResponse
	query := `
		SELECT u.email, p.full_name, p.phone, p.photo 
		FROM users u 
		JOIN profiles p ON u.id = p.user_id 
		WHERE u.id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&res.Email, &res.FullName, &res.Phone, &res.Photo)
	return res, err
}

func (r *userRepository) GetDashboard(ctx context.Context, userID int) (dto.DashboardResponse, error) {
	var res dto.DashboardResponse

	// 1. Ambil Saldo
	err := r.db.QueryRow(ctx, `SELECT balance FROM wallets WHERE user_id = $1`, userID).Scan(&res.Balance)
	if err != nil {
		return res, err
	}

	// 2. Ambil Pemasukan (IN) - Gunakan COALESCE agar jika null menjadi 0
	r.db.QueryRow(ctx, `SELECT SUM(amount) FROM transactions WHERE user_id = $1 AND flow_type = 'IN'`, userID).Scan(&res.Income)

	// 3. Ambil Pengeluaran (OUT)
	r.db.QueryRow(ctx, `SELECT SUM(amount) FROM transactions WHERE user_id = $1 AND flow_type = 'OUT'`, userID).Scan(&res.Expense)

	return res, nil
}
func (r *userRepository) GetUserByID(ctx context.Context, userID int) (model.User, error) {
	query := `SELECT id, email FROM users WHERE id = $1`

	var user model.User
	err := r.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Email)
	return user, err
}
func (r *userRepository) UpdateProfile(ctx context.Context, userID int, req dto.EditProfileRequest) error {
	query := `UPDATE profiles SET full_name = $1, phone = $2, photo = $3, updated_at = NOW() WHERE user_id = $4`
	_, err := r.db.Exec(ctx, query, req.FullName, req.Phone, req.Photo, userID)
	return err
}

func (r *userRepository) GetPasswordAndPin(ctx context.Context, userID int) (string, string, error) {
	var pass, pin string
	// Gunakan COALESCE untuk PIN berjaga-jaga jika masih NULL di database
	query := `SELECT password, COALESCE(pin, '') FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&pass, &pin)
	return pass, pin, err
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID int, newPassword string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`, newPassword, userID)
	return err
}

func (r *userRepository) UpdatePin(ctx context.Context, userID int, newPin string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET pin = $1, updated_at = NOW() WHERE id = $2`, newPin, userID)
	return err
}
