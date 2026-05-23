package pkg

import (
	"regexp"
)

// regex email untuk pengecekan format email
func IsValidEmail(email string) bool {
	regexPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(regexPattern)

	return re.MatchString(email)
}
