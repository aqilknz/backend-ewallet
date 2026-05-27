package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context) (*redis.Client, error) {
	_ = godotenv.Load()

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	username := os.Getenv("REDIS_USERNAME")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	addr := fmt.Sprintf("%s:%s", host, port)

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("gagal terhubung ke redis: %w", err)
	}

	log.Println("Redis Connected Successfully")
	return rdb, nil
}
