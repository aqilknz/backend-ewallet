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
	_ = godotenv.Load()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat connection pool: %w", err)
	}

	// cek koneksi benar-benar tembus ke database
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	log.Println("Database Connected Successfully")
	return pool, nil
}
