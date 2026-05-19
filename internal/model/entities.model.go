package model

import "time"

type User struct {
	ID        uint      `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Pin       string    `json:"-" db:"pin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Relasi (Gunakan omitempty agar tidak muncul di JSON jika datanya kosong/tidak di-join)
	Profile      Profile       `json:"profile,omitempty"`
	Wallet       Wallet        `json:"wallet,omitempty"`
	Topups       []Topup       `json:"topups,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

type Profile struct {
	ID        uint      `json:"id" db:"id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	FullName  string    `json:"full_name" db:"full_name"`
	Phone     string    `json:"phone" db:"phone"`
	Photo     string    `json:"photo" db:"photo"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Wallet struct {
	ID        uint      `json:"id" db:"id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	Balance   int       `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PaymentMethod struct {
	ID   uint   `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Topup struct {
	ID              uint      `json:"id" db:"id"`
	UserID          uint      `json:"user_id" db:"user_id"`
	PaymentMethodID uint      `json:"payment_method_id" db:"payment_method_id"`
	Amount          int       `json:"amount" db:"amount"`
	Fee             int       `json:"fee" db:"fee"`
	Status          string    `json:"status" db:"status"`
	Notes           string    `json:"notes" db:"notes"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`

	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
}

type Transfer struct {
	ID         uint      `json:"id" db:"id"`
	SenderID   uint      `json:"sender_id" db:"sender_id"`
	ReceiverID uint      `json:"receiver_id" db:"receiver_id"`
	Amount     int       `json:"amount" db:"amount"`
	Fee        int       `json:"fee" db:"fee"`
	Status     string    `json:"status" db:"status"`
	Notes      string    `json:"notes" db:"notes"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	Sender   User `json:"sender,omitempty"`
	Receiver User `json:"receiver,omitempty"`
}

type Transaction struct {
	ID              uint      `json:"id" db:"id"`
	UserID          uint      `json:"user_id" db:"user_id"`
	TransactionType string    `json:"transaction_type" db:"transaction_type"`
	FlowType        string    `json:"flow_type" db:"flow_type"` // contoh: IN / OUT
	Amount          int       `json:"amount" db:"amount"`
	ReferenceID     int       `json:"reference_id" db:"reference_id"`
	Description     string    `json:"description" db:"description"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
