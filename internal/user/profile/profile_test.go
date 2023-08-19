package profile_test

import (
	"github.com/google/uuid"
	"testing"
	"time"
	"user-service/internal/user/profile"
)

func TestNewProfile(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		dob       string
		email     string
		phone     string
		wantErr   bool
	}{
		{"Valid input", "John", "Doe", "2000-01-01", "john@example.com", "+123456789", false},
		{"Empty first name", "", "Doe", "2000-01-01", "john@example.com", "+123456789", false},
		{"Empty last name", "John", "", "2000-01-01", "john@example.com", "+123456789", false},
		{"Invalid dob format", "John", "Doe", "01-01-2000", "john@example.com", "+123456789", true},
		{"Empty email", "John", "Doe", "2000-01-01", "", "+123456789", false},
		{"Invalid email format", "John", "Doe", "2000-01-01", "johnexample.com", "+123456789", true},
		{"Empty phone", "John", "Doe", "2000-01-01", "john@example.com", "", false},
		{"Invalid phone format", "John", "Doe", "2000-01-01", "john@example.com", "1234", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := profile.NewProfile(tt.firstName, tt.lastName, tt.dob, tt.email, tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMapFromData(t *testing.T) {
	validUUID := uuid.New()
	tests := []struct {
		name      string
		profileID uuid.UUID
		firstName string
		lastName  string
		dob       time.Time
		email     string
		phone     string
		wantErr   bool
	}{
		{"Valid input", validUUID, "John", "Doe", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "john@example.com", "+123456789", false},
		{"Invalid email format", validUUID, "John", "Doe", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "john@.com", "+123456789", true},
		{"Invalid phone format", validUUID, "John", "Doe", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "john@example.com", "1234", true},
		{"Missing first name", validUUID, "", "Doe", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "john@example.com", "+123456789", false},
		{"Missing last name", validUUID, "John", "", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "john@example.com", "+123456789", false},
		{"Long first name", validUUID, "JohnJohnJohnJohnJohnJohnJohnJohn", "Doe", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "john@example.com", "+123456789", true},
		{"Long last name", validUUID, "John", "DoeDoeDoeDoeDoeDoeDoeDoe", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "john@example.com", "+123456789", true},
		{"Invalid UUID", uuid.Nil, "John", "Doe", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "john@example.com", "+123456789", true},
		{"All empty values", validUUID, "", "", time.Time{}, "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := profile.MapFromData(tt.profileID, tt.firstName, tt.lastName, tt.dob, tt.email, tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapFromData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
