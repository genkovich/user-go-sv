package profile_handler

import (
	"github.com/google/uuid"
	"user-service/internal/user/profile"
)

type Handler struct {
	Storage profile.Storage
}

func NewProfileHandler(storage profile.Storage) *Handler {
	return &Handler{
		Storage: storage,
	}
}

func (h *Handler) GetProfile(userId uuid.UUID) (*profile.Profile, error) {
	return h.Storage.GetByUserId(userId)
}

func (h *Handler) UpdateProfile(updateProfile UpdateProfileCommand) error {
	profileEntity, err := profile.NewProfile(
		updateProfile.FirstName,
		updateProfile.LastName,
		updateProfile.Dob,
		updateProfile.Email,
		updateProfile.Phone,
	)

	if err != nil {
		return err
	}

	err = h.Storage.Upsert(updateProfile.UserId, profileEntity)

	if err != nil {
		return err
	}

	return nil
}
