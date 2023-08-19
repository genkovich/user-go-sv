package user_handler

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"user-service/internal/user"
	"user-service/internal/user/token"
)

type Handler struct {
	Storage user.Storage
	JwtKey  []byte
}

func NewHandler(storage user.Storage, jwtSecret []byte) *Handler {
	return &Handler{
		Storage: storage,
		JwtKey:  jwtSecret,
	}
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

func (h *Handler) Authenticate(login string, password string) (*token.Token, error) {
	credentials, err := h.Storage.GetByLogin(login)

	if err != nil {
		return nil, err
	}

	if !credentials.IsCorrectPassword(password) {
		return nil, errors.New("invalid password")
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &token.Claims{
		Login: credentials.GetLogin(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(h.JwtKey)

	if err != nil {
		return nil, err
	}

	return &token.Token{JwtToken: tokenString}, nil
}
