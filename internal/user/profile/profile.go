package profile

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
	"user-service/internal/user/profile/field"
)

type Profile struct {
	id        uuid.UUID
	firstName field.FirstName
	lastName  field.LastName
	dob       time.Time
	email     field.Email
	phone     field.Phone
}

func NewProfile(firstName string, lastName string, dob string, email string, phone string) (*Profile, error) {
	var validateErrors []error

	var err error

	var firstNameField field.FirstName
	if firstName != "" {
		firstNameField, err = field.NewFirstName(firstName)
		if err != nil {
			validateErrors = append(validateErrors, err)
		}
	}

	var lastNameField field.LastName
	if lastName != "" {
		lastNameField, err = field.NewLastName(lastName)
		if err != nil {
			validateErrors = append(validateErrors, err)
		}
	}

	var dobField time.Time
	if dob != "" {
		dobField, err = time.Parse("2006-01-02", dob)
		if err != nil {
			validateErrors = append(validateErrors, err)
		}
	}

	var emailField field.Email
	if email != "" {
		emailField, err = field.NewEmail(email)
		if err != nil {
			validateErrors = append(validateErrors, err)
		}
	}

	var phoneField field.Phone
	if phone != "" {
		phoneField, err = field.NewPhone(phone)
		if err != nil {
			validateErrors = append(validateErrors, err)
		}
	}

	if len(validateErrors) > 0 {
		var resultError string
		for _, err := range validateErrors {
			resultError += err.Error() + "\n"
		}
		return nil, fmt.Errorf(resultError)
	}

	return &Profile{
		id:        uuid.New(),
		firstName: firstNameField,
		lastName:  lastNameField,
		dob:       dobField,
		email:     emailField,
		phone:     phoneField,
	}, nil
}

func (p *Profile) GetId() uuid.UUID {
	return p.id
}

func (p *Profile) GetFirstName() field.FirstName {
	return p.firstName
}

func (p *Profile) GetLastName() field.LastName {
	return p.lastName
}

func (p *Profile) GetDob() time.Time {
	return p.dob
}

func (p *Profile) GetEmail() field.Email {
	return p.email
}

func (p *Profile) GetPhone() field.Phone {
	return p.phone
}

func MapFromData(profileID uuid.UUID, firstName string, lastName string, dob time.Time, email string, phone string) (*Profile, error) {
	if (profileID == uuid.UUID{}) {
		return nil, fmt.Errorf("invalid data")
	}

	mappedProfile, err := NewProfile(firstName, lastName, dob.Format("2006-01-02"), email, phone)
	if err != nil {
		return nil, err
	}

	mappedProfile.id = profileID
	return mappedProfile, nil
}

func (p *Profile) MarshalJSON() ([]byte, error) {
	var dob *time.Time
	if !p.dob.IsZero() {
		dob = &p.dob
	}

	return json.Marshal(&struct {
		Id        uuid.UUID       `json:"id"`
		FirstName field.FirstName `json:"first_name,omitempty"`
		LastName  field.LastName  `json:"last_name,omitempty"`
		Dob       *time.Time      `json:"dob,omitempty"`
		Email     field.Email     `json:"email,omitempty"`
		Phone     field.Phone     `json:"phone,omitempty"`
	}{
		Id:        p.id,
		FirstName: p.firstName,
		LastName:  p.lastName,
		Dob:       dob,
		Email:     p.email,
		Phone:     p.phone,
	})
}
