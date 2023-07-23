package user_handler

import (
	"github.com/google/uuid"
	"user-service/internal/user"
)

type Handler struct {
	Storage user.Storage
}

func (h *Handler) GetList() []user.User {
	return h.Storage.GetList()
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
