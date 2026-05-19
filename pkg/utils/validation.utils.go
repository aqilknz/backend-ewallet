// internal/utils/validation.go
package utils

import (
	"regexp"
)

// IsValidEmail mengecek apakah string sesuai dengan standar format email yang ketat
func IsValidEmail(email string) bool {
	// Pola Regex:
	// 1. Harus diawali huruf/angka/karakter khusus yang diizinkan sebelum @
	// 2. Harus ada @
	// 3. Domain harus huruf/angka dengan ekstensi minimal 2 karakter (contoh: .com, .id, .co.id)
	regexPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Gunakan MustCompile agar regex dikompilasi sekali saja
	re := regexp.MustCompile(regexPattern)

	return re.MatchString(email)
}
