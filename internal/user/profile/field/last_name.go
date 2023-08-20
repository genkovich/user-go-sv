package field

import "fmt"

type LastName string

func NewLastName(name string) (LastName, error) {
	if len(name) > 20 || name == "" {
		return "", fmt.Errorf("invalid name, must be withn 20 charcters and non-empty")
	}

	return LastName(name), nil
}
