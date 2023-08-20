package field

import "fmt"

type FirstName string

func NewFirstName(name string) (FirstName, error) {
	if len(name) > 20 || name == "" {
		return "", fmt.Errorf("invalid name, must be withn 20 charcters and non-empty")
	}

	return FirstName(name), nil
}
