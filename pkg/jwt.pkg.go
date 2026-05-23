package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "fallback_secret_key"
	}

	return token.SignedString([]byte(secret))
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func VerifyToken(tokenString string) (int, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode enkripsi token tidak valid")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	// Ambil isinya (claims) jika token valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("data user_id tidak ditemukan di token")
		}
		return int(userIDFloat), nil
	}

	return 0, errors.New("token tidak valid")
}
