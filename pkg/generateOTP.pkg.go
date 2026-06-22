package pkg

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

func GenerateOTP() (string, error) {
	var otp string
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += strconv.Itoa(int(num.Int64()))
	}
	return otp, nil
}
