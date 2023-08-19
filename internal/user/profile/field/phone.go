package field

import (
	"fmt"
	"regexp"
	"strings"
)

type Phone struct {
	phone string
}

func NewPhone(phone string) (Phone, error) {
	phoneNormalized := normalize(phone)
	phoneNormalizedLen := len(phoneNormalized) - 1

	if phoneNormalizedLen > 20 || phoneNormalizedLen < 6 {
		return Phone{}, fmt.Errorf("invalid phone, must be less than 20 digits and more than 6 digits")
	}

	return Phone{phone: phone}, nil
}

func (p *Phone) String() string {
	return p.phone
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
