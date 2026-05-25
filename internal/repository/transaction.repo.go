package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/aqilknz/backend-ewallet/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	TopUp(ctx context.Context, userID int, req dto.TopUpRequest) (dto.TopUpResponse, error)
	Transfer(ctx context.Context, senderID int, receiverID int, req dto.TransferRequest) (dto.TransferResponse, error)
	GetUserIDByEmail(ctx context.Context, email string) (int, error)
	GetHistory(ctx context.Context, userID int, search string, limit int, offset int) ([]dto.TransactionHistoryItem, int, error)
	GetReport(ctx context.Context, userID int, param dto.TransactionReportFilterParam) ([]dto.TransactionReportItem, error)
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&id)
	return id, err
}

func (r *transactionRepository) TopUp(ctx context.Context, userID int, req dto.TopUpRequest) (dto.TopUpResponse, error) {
	var res dto.TopUpResponse
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return res, err
	}
	defer tx.Rollback(ctx)

	// Insert ke tabel transactions
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, amount, type, status, updated_at)
		VALUES ($1, $2, 'topup', 'success', NOW())
		RETURNING id, amount, status, created_at`,
		userID, req.Amount).Scan(&res.TransactionID, &res.Amount, &res.Status, &res.CreatedAt)
	if err != nil {
		return res, err
	}

	// Insert ke tabel topup_details
	_, err = tx.Exec(ctx, `
		INSERT INTO topup_details (transaction_id, payment_method_id, discount, tax, sub_total)
		VALUES ($1, $2, $3, $4, $5)`,
		res.TransactionID, req.PaymentMethodID, req.Discount, req.Tax, req.SubTotal)
	if err != nil {
		return res, err
	}

	// Update saldo wallet langsung ke balance
	_, err = tx.Exec(ctx, `UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2`, req.Amount, userID)
	if err != nil {
		return res, err
	}

	res.PaymentMethodID = req.PaymentMethodID
	res.Discount = req.Discount
	res.Tax = req.Tax
	res.SubTotal = req.SubTotal

	return res, tx.Commit(ctx)
}

func (r *transactionRepository) Transfer(ctx context.Context, senderID int, receiverID int, req dto.TransferRequest) (dto.TransferResponse, error) {
	var res dto.TransferResponse
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return res, err
	}
	defer tx.Rollback(ctx)

	// Cek saldo pengirim
	var currentBalance int
	err = tx.QueryRow(ctx, `SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE`, senderID).Scan(&currentBalance)
	if err != nil || currentBalance < req.Amount {
		return res, errors.New("saldo tidak mencukupi")
	}

	// Transaksi Pengirim (TRANSFER_OUT)
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, amount, type, status, updated_at)
		VALUES ($1, $2, 'transfer_out', 'success', NOW())
		RETURNING id, user_id, amount, status, created_at`,
		senderID, req.Amount).Scan(&res.TransactionID, &res.SenderID, &res.Amount, &res.Status, &res.CreatedAt)
	if err != nil {
		return res, err
	}

	// Detail untuk pengirim (countryparty nya itu penerima)
	_, err = tx.Exec(ctx, `
		INSERT INTO transfer_details (transaction_id, counterparty_id, notes) 
		VALUES ($1, $2, $3)`,
		res.TransactionID, receiverID, req.Notes)
	if err != nil {
		return res, err
	}

	//Transaksi Penerima (transfer_in)
	var receiverTxID int
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, amount, type, status, updated_at)
		VALUES ($1, $2, 'transfer_in', 'success', NOW())
		RETURNING id`,
		receiverID, req.Amount).Scan(&receiverTxID)
	if err != nil {
		return res, err
	}

	// Detail untuk penerima (countryparty nya itu pengirim)
	_, err = tx.Exec(ctx, `
		INSERT INTO transfer_details (transaction_id, counterparty_id, notes) 
		VALUES ($1, $2, $3)`,
		receiverTxID, senderID, req.Notes)
	if err != nil {
		return res, err
	}

	// Potong saldo pengirim dan tambah saldo penerima
	_, err = tx.Exec(ctx, `UPDATE wallets SET balance = balance - $1, updated_at = NOW() WHERE user_id = $2`, req.Amount, senderID)
	if err != nil {
		return res, err
	}
	_, err = tx.Exec(ctx, `UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2`, req.Amount, receiverID)
	if err != nil {
		return res, err
	}

	res.CounterpartyID = receiverID
	res.Notes = req.Notes

	return res, tx.Commit(ctx)
}

func (r *transactionRepository) GetHistory(ctx context.Context, userID int, search string, limit int, offset int) ([]dto.TransactionHistoryItem, int, error) {
	var histories []dto.TransactionHistoryItem = []dto.TransactionHistoryItem{}
	var totalRecords int
	searchParam := "%" + search + "%"

	baseQuery := `
		FROM transactions t
		LEFT JOIN transfer_details td ON t.id = td.transaction_id
		LEFT JOIN profiles p_counterparty ON td.counterparty_id = p_counterparty.user_id
		WHERE t.user_id = $1
	`

	countQuery := fmt.Sprintf(`
		SELECT COUNT(t.id) %s 
		AND (
			CASE 
				WHEN t.type = 'topup' THEN 'Topup Saldo'
				WHEN t.type = 'transfer_out' THEN CONCAT('Transfer keluar ke ', p_counterparty.full_name)
				WHEN t.type = 'transfer_in' THEN CONCAT('Transfer masuk dari ', p_counterparty.full_name)
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
			t.type::text,
			CASE 
				WHEN t.type = 'topup' THEN 'topup'
				WHEN t.type = 'transfer_out' THEN 'expense'
				WHEN t.type = 'transfer_in' THEN 'income'
			END as flow_type,
			CASE 
				WHEN t.type = 'topup' THEN 'Topup Saldo'
				WHEN t.type = 'transfer_out' THEN CONCAT('Transfer keluar ke ', p_counterparty.full_name)
				WHEN t.type = 'transfer_in' THEN CONCAT('Transfer masuk dari ', p_counterparty.full_name)
			END as description,
			t.created_at
		%s
		AND (
			CASE 
				WHEN t.type = 'topup' THEN 'Topup Saldo'
				WHEN t.type = 'transfer_out' THEN CONCAT('Transfer keluar ke ', p_counterparty.full_name)
				WHEN t.type = 'transfer_in' THEN CONCAT('Transfer masuk dari ', p_counterparty.full_name)
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
			SUM(CASE WHEN $2 IN ('all', 'income') AND t.type = 'transfer_in' THEN t.amount ELSE 0 END) AS total_income,
			SUM(CASE WHEN $2 IN ('all', 'expense') AND t.type = 'transfer_out' THEN t.amount ELSE 0 END) AS total_expense
		FROM transactions t
		WHERE t.user_id = $1
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
