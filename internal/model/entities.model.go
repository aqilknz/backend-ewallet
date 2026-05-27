package model

import "time"

type User struct {
	ID        uint      `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Pin       string    `json:"-" db:"pin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Profile      *Profile      `json:"profile,omitempty"`
	Wallet       *Wallet       `json:"wallet,omitempty"`
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

type Transaction struct {
	ID        uint      `json:"id" db:"id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	Amount    int       `json:"amount" db:"amount"`
	Type      string    `json:"type" db:"type"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	TopupDetail    *TopupDetail    `json:"topup_detail,omitempty"`
	TransferDetail *TransferDetail `json:"transfer_detail,omitempty"`
}

type TopupDetail struct {
	ID              uint `json:"id" db:"id"`
	TransactionID   uint `json:"transaction_id" db:"transaction_id"`
	PaymentMethodID uint `json:"payment_method_id" db:"payment_method_id"`
	Discount        int  `json:"discount" db:"discount"`
	Tax             int  `json:"tax" db:"tax"`
	SubTotal        int  `json:"sub_total" db:"sub_total"`

	PaymentMethod *PaymentMethod `json:"payment_method,omitempty"`
}

type TransferDetail struct {
	ID            uint   `json:"id" db:"id"`
	TransactionID uint   `json:"transaction_id" db:"transaction_id"`
	ReceiverID    uint   `json:"receiver_id" db:"receiver_id"`
	Notes         string `json:"notes" db:"notes"`

	Recipient *User `json:"recipient,omitempty"`
}
