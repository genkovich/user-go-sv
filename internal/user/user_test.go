package user_test

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
	"user-service/internal/user"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		passwordHash  string
		expectedError bool
	}{
		{"Valid creation #1", "testLogin1", "password1", false},
		{"Valid creation #2", "testLogin2", "password2", false},
		{"Empty login", "", "password", true},
		{"Empty password", "testLogin", "", true},
		{"Both empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := user.CreateUser(tt.login, tt.passwordHash)

			if tt.expectedError {
				assert.Nil(t, u)
				assert.NotNil(t, err)
			} else {
				assert.NotNil(t, u)
				assert.Nil(t, err)
				assert.Equal(t, tt.login, u.GetLogin())
				tempRole := u.GetRole()
				assert.Equal(t, "ROLE_USER", tempRole.GetRole())
			}
		})
	}
}

func TestMapFromData(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		login        string
		password     string
		role         string
		expectedUser *user.User
		expectedErr  bool
	}{
		{"Valid data #1", "550e8400-e29b-41d4-a716-446655440000", "testLogin1", "password1", "ROLE_USER", &user.User{}, false},
		{"Valid data #2", "550e8400-e29b-41d4-a716-446655440001", "testLogin2", "password2", "ROLE_USER", &user.User{}, false},
		{"Valid data #3", "550e8400-e29b-41d4-a716-446655440002", "testLogin3", "password3", "ROLE_ADMIN", &user.User{}, false},
		{"Valid data #4", "550e8400-e29b-41d4-a716-446655440003", "testLogin4", "password4", "ROLE_USER", &user.User{}, false},
		{"Valid data #5", "550e8400-e29b-41d4-a716-446655440004", "testLogin5", "password5", "ROLE_ADMIN", &user.User{}, false},

		{"Invalid UUID", "invalid-uuid", "testLogin", "password", "ROLE_USER", nil, true},
		{"Empty UUID", "", "testLogin", "password", "User", nil, true},
		{"Invalid Role", "550e8400-e29b-41d4-a716-446655440005", "testLogin", "password", "InvalidRole", nil, true},
		{"Empty Login", "550e8400-e29b-41d4-a716-446655440006", "", "password", "ROLE_USER", nil, true},
		{"Empty Password", "550e8400-e29b-41d4-a716-446655440007", "testLogin", "", "ROLE_USER", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userResult, err := user.MapFromData(tt.id, tt.login, tt.password, tt.role)

			if tt.expectedErr {
				assert.NotNil(t, err)
				assert.Nil(t, userResult)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, userResult)
				assert.Equal(t, tt.id, userResult.GetId().String())
				assert.Equal(t, tt.login, userResult.GetLogin())
				tempRole := userResult.GetRole()
				assert.Equal(t, tt.role, tempRole.GetRole())

			}
		})
	}
}

func TestChangeRole(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		newRole  string
		expected string
		err      bool
	}{
		{"Change from User to Admin", "ROLE_USER", "ROLE_ADMIN", "ROLE_ADMIN", false},
		{"Change from Admin to User", "ROLE_ADMIN", "ROLE_USER", "ROLE_USER", false},
		{"Invalid Role", "ROLE_USER", "InvalidRole", "ROLE_USER", true},
		{"Empty Role", "ROLE_ADMIN", "", "ROLE_ADMIN", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRole, _ := user.TryFrom(tt.initial)
			u, _ := user.CreateUser("login", "password")

			u.ChangeRole(*userRole)

			newUserRole, _ := user.TryFrom(tt.newRole)
			if newUserRole == nil {
				tempRole := u.GetRole()
				assert.Equal(t, tt.expected, tempRole.GetRole())
			} else {
				u.ChangeRole(*newUserRole)

				tempRole := u.GetRole()
				assert.Equal(t, tt.expected, tempRole.GetRole())
			}
		})
	}
}

func TestIsCorrectPassword(t *testing.T) {
	password := "testPassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}

	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Correct password", password, true},
		{"Incorrect password #1", "wrongPassword", false},
		{"Incorrect password #2", "AnotherWrongPassword", false},
		{"Empty password", "", false},
		{"Special characters", "#$%^&*()", false},
		{"Long incorrect password", "ThisIsAVeryLongPasswordWhichIsIncorrect", false},
	}

	u, _ := user.CreateUser("login", string(hashedPassword))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := u.IsCorrectPassword(tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}
