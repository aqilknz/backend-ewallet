package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// connect ke database
	godotenv.Load(".env")
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.Fatal("Gagal koneksi:", err)
	}
	defer db.Close()

	tx, _ := db.Begin()
	defer tx.Rollback()

	// insert payment methods
	paymentMethods := []string{"Bank Rakyat Indonesia", "DANA", "Bank Central Asia", "GoPay", "OVO"}
	for _, pm := range paymentMethods {
		tx.Exec(`INSERT INTO payment_methods (name) VALUES ($1) ON CONFLICT DO NOTHING`, pm)
	}

	// Generate 10 Users & Relasinya (Profiles & Wallets)
	hashPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashPin, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	// nyimpan user yang nanti dibuat
	var userIDs []int

	for i := 1; i <= 10; i++ {
		var userID int

		// Insert User dan ambil ID-nya
		tx.QueryRow(`INSERT INTO users (email, password, pin) VALUES ($1, $2, $3) RETURNING id`,
			fmt.Sprintf("user%d@mail.com", i), string(hashPass), string(hashPin),
		).Scan(&userID)

		// simpan ID untuk transaksi nanti
		userIDs = append(userIDs, userID)

		// insert di profile dan wallets
		tx.Exec(`INSERT INTO profiles (user_id, full_name, phone, photo) VALUES ($1, $2, $3, $4)`,
			userID, fmt.Sprintf("Pengguna Ke-%d", i), fmt.Sprintf("0812345678%d", i), fmt.Sprintf("https://i.pravatar.cc/150?u=%d", userID))

		tx.Exec(`INSERT INTO wallets (user_id, balance) VALUES ($1, $2)`, userID, i*1000000)
	}

	// simulasi transaksi topup dan transfer
	var topupID, tfID int

	// topup
	tx.QueryRow(`INSERT INTO transactions (user_id, amount, type, status) VALUES ($1, 500000, 'topup', 'success') RETURNING id`, userIDs[0]).Scan(&topupID)
	tx.Exec(`INSERT INTO topup_details (transaction_id, payment_method_id, discount, tax, sub_total) VALUES ($1, 1, 0, 2000, 502000)`, topupID)

	tx.QueryRow(`INSERT INTO transactions (user_id, amount, type, status) VALUES ($1, 150000, 'transfer_out', 'success') RETURNING id`, userIDs[0]).Scan(&tfID)
	tx.Exec(`INSERT INTO transfer_details (transaction_id, receiver_id, notes) VALUES ($1, $2, 'Makan siang')`, tfID, userIDs[1])

	tx.Exec(`INSERT INTO transactions (user_id, amount, type, status) VALUES ($1, 150000, 'transfer_in', 'success')`, userIDs[1])

	// nyimpan perubahan
	if err := tx.Commit(); err != nil {
		log.Fatal("Gagal menyimpan data:", err)
	}

	fmt.Println("✅ Seeding Sukses! Silakan login di React dengan 'user1@mail.com' / 'password123'")
}
