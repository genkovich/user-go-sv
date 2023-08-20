package field

import (
	"fmt"
	"regexp"
)

type Email string

func NewEmail(email string) (Email, error) {
	if len(email) > 200 || email == "" {
		return "", fmt.Errorf("invalid email, must be withn 200 charcters and non-empty")
	}

	reg := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)
	if !reg.MatchString(email) {
		return "", fmt.Errorf("invalid email, must be a valid email address")
	}

	return Email(email), nil
}
