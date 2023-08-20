package profile_handler

import "github.com/google/uuid"

type UpdateProfileCommand struct {
	UserId    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Dob       string    `json:"dob" validate:"omitempty,datetime=2006-01-02"`
	Email     string    `json:"email" validate:"omitempty,email,min=6,max=200"`
	Phone     string    `json:"phone" validate:"omitempty,min=6,max=50"`
}
