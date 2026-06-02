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
	EditProfile(ctx context.Context, userID int, fullname *string, phone *string, pictureURL *string) error
	GetPasswordAndPin(ctx context.Context, userID int) (password string, pin string, err error)
	UpdatePassword(ctx context.Context, userID int, newPassword string) error
	UpdatePin(ctx context.Context, userID int, newPin string) error
	GetUserByID(ctx context.Context, userID int) (model.User, error)
	FindReceivers(ctx context.Context, userID int, search string, limit int, offset int) ([]dto.ReceiverResponse, int, error)
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

	// amvil saldo wallet
	err := r.db.QueryRow(ctx, `SELECT balance FROM wallets WHERE user_id = $1`, userID).Scan(&res.Balance)
	if err != nil {
		return res, err
	}

	// income (transfer in)
	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(t.amount), 0) 
		FROM transactions t
		JOIN transfer_details td ON t.id = td.transaction_id
		WHERE t.type = 'transfer_out' AND td.receiver_id = $1
	`, userID).Scan(&res.Income)
	if err != nil {
		return res, err
	}

	// transfer out
	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE user_id = $1 AND type = 'transfer_out'
	`, userID).Scan(&res.Expense)
	if err != nil {
		return res, err
	}

	return res, nil
}
func (r *userRepository) GetUserByID(ctx context.Context, userID int) (model.User, error) {
	query := `SELECT id, email FROM users WHERE id = $1`

	var user model.User
	err := r.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Email)
	return user, err
}
func (r *userRepository) EditProfile(ctx context.Context, userID int, fullname *string, phone *string, pictureURL *string) error {
	query := `
		UPDATE profiles 
		SET full_name = COALESCE($1, full_name), 
		    phone = COALESCE($2, phone), 
		    photo = COALESCE($3, photo), 
		    updated_at = NOW() 
		WHERE user_id = $4`

	_, err := r.db.Exec(ctx, query, fullname, phone, pictureURL, userID)
	return err
}

func (r *userRepository) GetPasswordAndPin(ctx context.Context, userID int) (string, string, error) {
	var pass, pin string
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

func (r *userRepository) FindReceivers(ctx context.Context, userID int, search string, limit int, offset int) ([]dto.ReceiverResponse, int, error) {
	var receivers []dto.ReceiverResponse = []dto.ReceiverResponse{}
	var totalRecords int
	searchParam := "%" + search + "%"

	countQuery := `
		SELECT COUNT(u.id)
		FROM users u
		JOIN profiles p ON u.id = p.user_id
		WHERE u.id != $1 AND (p.phone ILIKE $2 OR p.full_name ILIKE $2)`

	err := r.db.QueryRow(ctx, countQuery, userID, searchParam).Scan(&totalRecords)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT u.id, p.full_name, u.email, p.phone, p.photo
		FROM users u
		JOIN profiles p ON u.id = p.user_id
		WHERE u.id != $1 AND (p.phone ILIKE $2 OR p.full_name ILIKE $2)
		ORDER BY p.full_name ASC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(ctx, dataQuery, userID, searchParam, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var rec dto.ReceiverResponse
		if err := rows.Scan(&rec.ID, &rec.FullName, &rec.Email, &rec.Phone, &rec.Photo); err != nil {
			return nil, 0, err
		}
		receivers = append(receivers, rec)
	}

	return receivers, totalRecords, nil
}
