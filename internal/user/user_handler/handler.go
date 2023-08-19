package user_handler

import (
	"errors"
	"github.com/google/uuid"
	"user-service/internal/user"
	"user-service/internal/user/token"
)

type Handler struct {
	storage        user.Storage
	tokenGenerator token.Generator
}

func NewHandler(storage user.Storage, tokenGenerator token.Generator) *Handler {
	return &Handler{
		storage:        storage,
		tokenGenerator: tokenGenerator,
	}
}

func (h *Handler) GetList(limit int, offset int) ([]user.User, error) {
	return h.storage.GetList(limit, offset)
}

func (h *Handler) Create(createUser CreateUserCommand) uuid.UUID {
	userEntity := *user.CreateUser(createUser.Login, createUser.Password)
	h.storage.Add(userEntity)

	return userEntity.GetId()
}

func (h *Handler) Delete(userId string) error {
	h.storage.Remove(userId)

	return nil
}

func (h *Handler) Authenticate(login string, password string) (*token.Token, error) {
	credentials, err := h.storage.GetByLogin(login)

	if err != nil {
		return nil, err
	}

	if !credentials.IsCorrectPassword(password) {
		return nil, errors.New("invalid password")
	}

	role := credentials.GetRole()
	return h.tokenGenerator.Generate(credentials.GetLogin(), role.GetRole())
}
