package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aqilknz/backend-ewallet/internal/config"
	"github.com/aqilknz/backend-ewallet/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title						Backend Ewallet
// @version						1.0
// @description					Backend created by Aqil using Gin

// @license.name				MIT

// @host						localhost:9000
// @BasePath					/ewallet

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Bearer token used for authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env. \ncause: %s", err.Error())
	}
	// inisialisasi
	app := gin.Default()

	// connect ke db
	db, err := config.ConnectDB(context.Background())
	if err != nil {
		log.Printf("DB connection error. \ncause: %s", err.Error())
	}
	defer db.Close()
	log.Println("DB Connected")

	// connect ke redis
	redis, err := config.ConnectRedis(context.Background())
	if err != nil {
		log.Printf("Redis error: %v", err)
	}
	defer redis.Close()
	log.Println("Redis Connected")

	// install router
	router.InitRouter(app, db, redis)

	// run
	addr := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	if err := app.Run(addr); err != nil {
		log.Printf("Server gagal berjalan: %v", err)
	}
}
