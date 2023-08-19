package field

import "fmt"

type LastName struct {
	lastName string
}

func NewLastName(name string) (LastName, error) {
	if len(name) > 20 || name == "" {
		return LastName{}, fmt.Errorf("invalid name, must be withn 20 charcters and non-empty")
	}

	return LastName{lastName: name}, nil
}

func (l *LastName) String() string {
	return l.lastName
}
