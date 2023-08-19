package profile

import "github.com/google/uuid"

type Storage interface {
	GetByUserId(userId uuid.UUID) (*Profile, error)
	Upsert(userId uuid.UUID, profile *Profile) error
}
