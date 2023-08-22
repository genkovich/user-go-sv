package user

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id           uuid.UUID
	login        string
	passwordHash string
	role         Role
	email        string
	privateUser  *User
}

func CreateUser(login string, passwordHash string) (*User, error) {
	if login == "" || passwordHash == "" {
		return nil, fmt.Errorf("invalid data")
	}

	return &User{
		id:           uuid.New(),
		login:        login,
		passwordHash: passwordHash,
		role:         *CreateUserRole(),
	}, nil
}

func NewUser(login string, hashedPassword string, email string) *User {
	user := &User{
		id:           uuid.New(),
		login:        login,
		passwordHash: hashedPassword,
		role:         *CreateUserRole(),
		email:        email,
	}

	return user
}

func MapFromData(id string, login string, passwordHash string, role string) (*User, error) {
	if id == "" || login == "" || passwordHash == "" || role == "" {
		return nil, fmt.Errorf("invalid data")
	}
	parse, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	userRole, err := TryFrom(role)
	if err != nil {
		return nil, err
	}

	return &User{
		id:           parse,
		login:        login,
		passwordHash: passwordHash,
		role:         *userRole,
	}, nil
}

func (u *User) GetId() uuid.UUID {
	return u.id
}

func (u *User) ChangeRole(r Role) {
	u.role = r
}

func (u *User) GetRole() Role {
	return u.role
}

func (u *User) GetLogin() string {
	return u.login
}

func (u *User) GetHashPassword() string {
	return u.passwordHash
}

func (u *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id    uuid.UUID `json:"id"`
		Login string    `json:"login"`
		Role  string    `json:"role"`
	}{
		Id:    u.id,
		Login: u.login,
		Role:  u.role.GetRole(),
	})
}
