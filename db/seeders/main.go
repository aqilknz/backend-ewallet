package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/aqilknz/backend-ewallet/pkg"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Seeder ini membuat:
//   - 5 payment methods
//   - 10 user dummy (email user1@mail.com .. user10@mail.com, password "pass1234")
//   - profile + wallet untuk tiap user
//   - 20 transaksi dummy (campuran topup & transfer), lengkap dengan detail relasinya
//
// Password & PIN di-hash menggunakan argon2id (pkg.HashData), supaya konsisten
// dengan pkg.VerifyHash yang dipakai di internal/service/auth.service.go & user.service.go.
func main() {
	godotenv.Load(".env")

	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database belum siap menerima koneksi:", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Gagal memulai transaction:", err)
	}
	defer tx.Rollback()

	// ── Payment Methods ──────────────────────────────────────────────
	paymentMethods := []string{"Bank Rakyat Indonesia", "DANA", "Bank Central Asia", "GoPay", "OVO"}
	for _, pm := range paymentMethods {
		if _, err := tx.Exec(`INSERT INTO payment_methods (name) VALUES ($1) ON CONFLICT DO NOTHING`, pm); err != nil {
			log.Fatal("Gagal seed payment_methods:", err)
		}
	}

	// ── Users + Profiles + Wallets ───────────────────────────────────
	hashedPassword, err := pkg.HashData("pass1234")
	if err != nil {
		log.Fatal("Gagal hash password:", err)
	}
	hashedPin, err := pkg.HashData("123456")
	if err != nil {
		log.Fatal("Gagal hash pin:", err)
	}

	const totalUsers = 10
	userIDs := make([]int, 0, totalUsers)
	startingBalance := make(map[int]int)

	for i := 1; i <= totalUsers; i++ {
		var userID int

		err := tx.QueryRow(
			`INSERT INTO users (email, password, pin, created_at, updated_at) 
			 VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
			fmt.Sprintf("user%d@mail.com", i), hashedPassword, hashedPin,
		).Scan(&userID)
		if err != nil {
			log.Fatal("Gagal insert user:", err)
		}
		userIDs = append(userIDs, userID)

		// Avatar pakai pravatar.cc - URL gambar publik yang valid & konsisten
		photoURL := fmt.Sprintf("https://i.pravatar.cc/150?img=%d", i)

		if _, err := tx.Exec(
			`INSERT INTO profiles (user_id, full_name, phone, photo, created_at, updated_at) 
			 VALUES ($1, $2, $3, $4, NOW(), NOW())`,
			userID, fmt.Sprintf("Pengguna Ke-%d", i), fmt.Sprintf("0812345678%02d", i), photoURL,
		); err != nil {
			log.Fatal("Gagal insert profile:", err)
		}

		initialBalance := 1000000 * i // user1 = 1jt, user2 = 2jt, dst (variasi saldo awal)
		if _, err := tx.Exec(
			`INSERT INTO wallets (user_id, balance, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())`,
			userID, initialBalance,
		); err != nil {
			log.Fatal("Gagal insert wallet:", err)
		}
		startingBalance[userID] = initialBalance
	}

	// ── 20 Transaksi Dummy ────────────────────────────────────────────
	// Distribusi: 8 topup, 12 transfer (yang otomatis menghasilkan
	// pasangan transfer_out + transfer_in agar history/report query konsisten).
	notesPool := []string{
		"Bayar makan siang", "Patungan kado", "Bayar listrik", "Titip jajan",
		"Split bill", "Bayar parkir", "Ganti ongkir", "Bayar utang",
	}

	transactionCount := 0
	const totalTransactions = 20
	const topupCount = 8

	for transactionCount < topupCount {
		userID := userIDs[rand.Intn(totalUsers)]
		paymentMethodID := rand.Intn(len(paymentMethods)) + 1
		amount := (rand.Intn(20) + 1) * 50000 // kelipatan 50rb, 50rb - 1jt
		tax := 4000
		discount := 0
		subTotal := amount + tax - discount

		var txID int
		err := tx.QueryRow(
			`INSERT INTO transactions (user_id, amount, type, status, created_at, updated_at) 
			 VALUES ($1, $2, 'topup', 'success', NOW() - (random() * INTERVAL '30 days'), NOW()) 
			 RETURNING id`,
			userID, amount,
		).Scan(&txID)
		if err != nil {
			log.Fatal("Gagal insert transaksi topup:", err)
		}

		if _, err := tx.Exec(
			`INSERT INTO topup_details (transaction_id, payment_method_id, discount, tax, sub_total) 
			 VALUES ($1, $2, $3, $4, $5)`,
			txID, paymentMethodID, discount, tax, subTotal,
		); err != nil {
			log.Fatal("Gagal insert topup_details:", err)
		}

		if _, err := tx.Exec(`UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2`, amount, userID); err != nil {
			log.Fatal("Gagal update wallet (topup):", err)
		}

		transactionCount++
	}

	for transactionCount < totalTransactions {
		senderID := userIDs[rand.Intn(totalUsers)]
		receiverID := userIDs[rand.Intn(totalUsers)]
		for receiverID == senderID {
			receiverID = userIDs[rand.Intn(totalUsers)]
		}
		amount := (rand.Intn(10) + 1) * 25000 // kelipatan 25rb, 25rb - 250rb
		notes := notesPool[rand.Intn(len(notesPool))]

		var txID int
		err := tx.QueryRow(
			`INSERT INTO transactions (user_id, amount, type, status, created_at, updated_at) 
			 VALUES ($1, $2, 'transfer_out', 'success', NOW() - (random() * INTERVAL '30 days'), NOW()) 
			 RETURNING id`,
			senderID, amount,
		).Scan(&txID)
		if err != nil {
			log.Fatal("Gagal insert transaksi transfer:", err)
		}

		if _, err := tx.Exec(
			`INSERT INTO transfer_details (transaction_id, receiver_id, notes) VALUES ($1, $2, $3)`,
			txID, receiverID, notes,
		); err != nil {
			log.Fatal("Gagal insert transfer_details:", err)
		}

		if _, err := tx.Exec(`UPDATE wallets SET balance = balance - $1, updated_at = NOW() WHERE user_id = $2`, amount, senderID); err != nil {
			log.Fatal("Gagal update wallet (sender):", err)
		}
		if _, err := tx.Exec(`UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2`, amount, receiverID); err != nil {
			log.Fatal("Gagal update wallet (receiver):", err)
		}

		transactionCount++
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("Gagal menyimpan seed data:", err)
	}

	fmt.Println("Seeding selesai.")
	fmt.Println("- 10 user dummy dibuat (user1@mail.com .. user10@mail.com)")
	fmt.Println("- Password semua user: pass1234")
	fmt.Println("- PIN semua user: 123456")
	fmt.Println("- 20 transaksi dummy dibuat (8 topup, 12 transfer)")
}
