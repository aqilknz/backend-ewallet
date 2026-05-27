package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	CreateTopUp(ctx context.Context, userID int, req dto.TopUpRequest, tax, discount, subTotal int) (dto.TopUpResponse, error)
	CreateTransfer(ctx context.Context, senderID int, receiverID int, amount int, notes string) (dto.TransferResponse, error)
	GetSenderPin(ctx context.Context, userID int) (string, error)
	GetHistory(ctx context.Context, userID int, search string, limit int, offset int) ([]dto.TransactionHistoryItem, int, error)
	GetReport(ctx context.Context, userID int, param dto.TransactionReportFilterParam) ([]dto.TransactionReportItem, error)
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (tr *transactionRepository) GetSenderPin(ctx context.Context, userID int) (string, error) {
	var pin string
	err := tr.db.QueryRow(ctx, "SELECT COALESCE(pin, '') FROM users WHERE id = $1", userID).Scan(&pin)
	return pin, err
}

func (tr *transactionRepository) CreateTopUp(ctx context.Context, userID int, req dto.TopUpRequest, tax, discount, subTotal int) (dto.TopUpResponse, error) {
	var res dto.TopUpResponse

	tx, err := tr.db.Begin(ctx)
	if err != nil {
		return res, err
	}
	defer tx.Rollback(ctx)

	sqlBase := `
		INSERT INTO transactions (user_id, amount, type, status, created_at, updated_at) 
		VALUES ($1, $2, 'topup', 'success', NOW(), NOW()) 
		RETURNING id, created_at
	`
	err = tx.QueryRow(ctx, sqlBase, userID, req.Amount).Scan(&res.TransactionID, &res.CreatedAt)
	if err != nil {
		return res, err
	}

	sqlDetail := `
		INSERT INTO topup_details (transaction_id, payment_method_id, discount, tax, sub_total) 
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, sqlDetail, res.TransactionID, req.PaymentMethodID, discount, tax, subTotal)
	if err != nil {
		return res, err
	}

	sqlWallet := `UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2`
	_, err = tx.Exec(ctx, sqlWallet, req.Amount, userID)
	if err != nil {
		return res, err
	}

	if err := tx.Commit(ctx); err != nil {
		return res, err
	}

	res.Amount = req.Amount
	res.PaymentMethodID = req.PaymentMethodID
	res.Tax = tax
	res.Discount = discount
	res.SubTotal = subTotal
	res.Status = "success"

	return res, nil
}

func (tr *transactionRepository) CreateTransfer(ctx context.Context, senderID int, receiverID int, amount int, notes string) (dto.TransferResponse, error) {
	var res dto.TransferResponse

	tx, err := tr.db.Begin(ctx)
	if err != nil {
		return res, err
	}
	defer tx.Rollback(ctx)

	sqlDeduct := `UPDATE wallets SET balance = balance - $1, updated_at = NOW() WHERE user_id = $2 AND balance >= $1`
	cmdTag, err := tx.Exec(ctx, sqlDeduct, amount, senderID)
	if err != nil {
		return res, err
	}
	if cmdTag.RowsAffected() == 0 {
		return res, errors.New("saldo tidak mencukupi")
	}

	sqlAdd := `UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2`
	cmdTagReceiver, err := tx.Exec(ctx, sqlAdd, amount, receiverID)
	if err != nil {
		return res, err
	}
	if cmdTagReceiver.RowsAffected() == 0 {
		return res, errors.New("user penerima tidak ditemukan")
	}

	sqlBase := `
		INSERT INTO transactions (user_id, amount, type, status, created_at, updated_at) 
		VALUES ($1, $2, 'transfer_out', 'success', NOW(), NOW()) 
		RETURNING id, created_at
	`
	err = tx.QueryRow(ctx, sqlBase, senderID, amount).Scan(&res.TransactionID, &res.CreatedAt)
	if err != nil {
		return res, err
	}

	sqlDetail := `
		INSERT INTO transfer_details (transaction_id, receiver_id, notes) 
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, sqlDetail, res.TransactionID, receiverID, notes)
	if err != nil {
		return res, err
	}

	if err := tx.Commit(ctx); err != nil {
		return res, err
	}

	res.SenderID = senderID
	res.ReceiverID = receiverID
	res.Amount = amount
	res.Notes = notes
	res.Status = "success"

	return res, nil
}

func (r *transactionRepository) GetHistory(ctx context.Context, userID int, search string, limit int, offset int) ([]dto.TransactionHistoryItem, int, error) {
	var histories []dto.TransactionHistoryItem = []dto.TransactionHistoryItem{}
	var totalRecords int
	searchParam := "%" + search + "%"

	baseQuery := `
		FROM transactions t
		LEFT JOIN transfer_details td ON t.id = td.transaction_id
		LEFT JOIN profiles p_sender ON t.user_id = p_sender.user_id
		LEFT JOIN profiles p_receiver ON td.receiver_id = p_receiver.user_id
		WHERE (t.user_id = $1 OR td.receiver_id = $1)
	`

	countQuery := fmt.Sprintf(`
		SELECT COUNT(t.id) %s 
		AND (
			CASE 
				WHEN t.type = 'topup' THEN 'Topup Saldo'
				WHEN t.type = 'transfer_out' AND t.user_id = $1 THEN CONCAT('Transfer keluar ke ', p_receiver.full_name)
				WHEN t.type = 'transfer_out' AND td.receiver_id = $1 THEN CONCAT('Transfer masuk dari ', p_sender.full_name)
			END ILIKE $2
		)`, baseQuery)

	err := r.db.QueryRow(ctx, countQuery, userID, searchParam).Scan(&totalRecords)
	if err != nil {
		return nil, 0, err
	}

	mainQuery := fmt.Sprintf(`
		SELECT 
			t.id, 
			t.amount, 
			CASE 
				WHEN t.type = 'topup' THEN 'topup'
				WHEN t.type = 'transfer_out' AND t.user_id = $1 THEN 'transfer_out'
				WHEN t.type = 'transfer_out' AND td.receiver_id = $1 THEN 'transfer_in'
			END as transaction_type,
			CASE 
				WHEN t.type = 'topup' THEN 'topup'
				WHEN t.type = 'transfer_out' AND t.user_id = $1 THEN 'expense'
				WHEN t.type = 'transfer_out' AND td.receiver_id = $1 THEN 'income'
			END as flow_type,
			CASE 
				WHEN t.type = 'topup' THEN 'Topup Saldo'
				WHEN t.type = 'transfer_out' AND t.user_id = $1 THEN CONCAT('Transfer keluar ke ', p_receiver.full_name)
				WHEN t.type = 'transfer_out' AND td.receiver_id = $1 THEN CONCAT('Transfer masuk dari ', p_sender.full_name)
			END as description,
			t.created_at
		%s
		AND (
			CASE 
				WHEN t.type = 'topup' THEN 'Topup Saldo'
				WHEN t.type = 'transfer_out' AND t.user_id = $1 THEN CONCAT('Transfer keluar ke ', p_receiver.full_name)
				WHEN t.type = 'transfer_out' AND td.receiver_id = $1 THEN CONCAT('Transfer masuk dari ', p_sender.full_name)
			END ILIKE $2
		)
		ORDER BY t.created_at DESC LIMIT $3 OFFSET $4`, baseQuery)

	rows, err := r.db.Query(ctx, mainQuery, userID, searchParam, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var item dto.TransactionHistoryItem
		err := rows.Scan(&item.ID, &item.Amount, &item.TransactionType, &item.FlowType, &item.Description, &item.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		histories = append(histories, item)
	}

	return histories, totalRecords, nil
}

func (r *transactionRepository) GetReport(ctx context.Context, userID int, param dto.TransactionReportFilterParam) ([]dto.TransactionReportItem, error) {
	var reports []dto.TransactionReportItem = []dto.TransactionReportItem{}

	query := `
		SELECT 
			TO_CHAR(DATE(t.created_at), 'YYYY-MM-DD') AS date,
			SUM(CASE WHEN $2 IN ('all', 'income') AND t.type = 'transfer_out' AND td.receiver_id = $1 THEN t.amount ELSE 0 END) AS total_income,
			SUM(CASE WHEN $2 IN ('all', 'expense') AND t.type = 'transfer_out' AND t.user_id = $1 THEN t.amount ELSE 0 END) AS total_expense
		FROM transactions t
		LEFT JOIN transfer_details td ON t.id = td.transaction_id
		WHERE (t.user_id = $1 OR td.receiver_id = $1)
	`

	args := []any{userID, param.Type}
	argCount := 2

	if param.StartDate != "" {
		argCount++
		query += fmt.Sprintf(" AND DATE(t.created_at) >= $%d", argCount)
		args = append(args, param.StartDate)
	}
	if param.EndDate != "" {
		argCount++
		query += fmt.Sprintf(" AND DATE(t.created_at) <= $%d", argCount)
		args = append(args, param.EndDate)
	}

	query += ` GROUP BY DATE(t.created_at) ORDER BY DATE(t.created_at) ASC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item dto.TransactionReportItem
		if err := rows.Scan(&item.Date, &item.TotalIncome, &item.TotalExpense); err != nil {
			return nil, err
		}
		reports = append(reports, item)
	}

	return reports, nil
}
