package user_handler

import (
	"github.com/google/uuid"
	"user-service/internal/user"
)

type Handler struct {
	Storage user.Storage
}

func NewHandler(storage user.Storage) *Handler {
	return &Handler{Storage: storage}
}

func (h *Handler) GetList(limit int, offset int) ([]user.User, error) {
	return h.Storage.GetList(limit, offset)
}

func (h *Handler) Create(createUser CreateUserCommand) uuid.UUID {
	userEntity := *user.CreateUser(createUser.Login, createUser.Password)
	h.Storage.Add(userEntity)

	return userEntity.GetId()
}

func (h *Handler) Delete(userId string) error {
	h.Storage.Remove(userId)

	return nil
}
