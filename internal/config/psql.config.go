package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	// 1. Load .env (Kita abaikan error-nya karena di Production/Docker,
	// variabel biasanya di-inject langsung tanpa file .env)
	_ = godotenv.Load()

	// 2. Susun DSN (Data Source Name) langsung dari pemanggilan Getenv
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 3. Buat koneksi pool (pgxpool.New otomatis melakukan ParseConfig di belakang layar!)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat connection pool: %w", err)
	}

	// 4. Pastikan koneksi benar-benar tembus ke database
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	log.Println("Database Connected Successfully")
	return pool, nil
}
