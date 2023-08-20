package field

import (
	"fmt"
	"regexp"
	"strings"
)

type Phone string

func NewPhone(phone string) (Phone, error) {
	phoneNormalized := normalize(phone)
	phoneNormalizedLen := len(phoneNormalized) - 1

	if phoneNormalizedLen > 20 || phoneNormalizedLen < 6 {
		return "", fmt.Errorf("invalid phone, must be less than 20 digits and more than 6 digits")
	}

	return Phone(phone), nil
}

func normalize(phone string) string {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(phone, -1)
	digits := strings.Join(matches, "")
	if len(digits) == 0 {
		return ""
	}

	return fmt.Sprintf("+%s", strings.Join(matches, ""))
}
