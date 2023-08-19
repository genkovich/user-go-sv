package field

import "fmt"

type FirstName struct {
	firstName string
}

func NewFirstName(name string) (FirstName, error) {
	if len(name) > 20 || name == "" {
		return FirstName{}, fmt.Errorf("invalid name, must be withn 20 charcters and non-empty")
	}

	return FirstName{firstName: name}, nil
}

func (n FirstName) String() string {
	return n.firstName
}
