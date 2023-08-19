package user

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"user-service/internal/user"
	"user-service/internal/user/user_handler"
	"user-service/pkg/response"
)

type Controller struct {
	handler *user_handler.Handler
	log     *zap.Logger
}

type RegisterUserCommand struct {
	Login           string `json:"login"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Email           string `json:"email"`
}

func NewUserController(handler *user_handler.Handler, log *zap.Logger) *Controller {
	return &Controller{
		handler: handler,
		log:     log,
	}
}

func (uc *Controller) RegisterRoutes(r *chi.Mux) {
	r.Get("/users", uc.GetUserList)
	r.Post("/users", uc.CreateUser)
	r.Delete("/users/{userId}", uc.DeleteUser)
	r.Post("/users/auth", uc.Authenticate)
	r.Post("/users/register", uc.RegisterUser)
}

func (uc *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var createUser RegisterUserCommand

	err := json.NewDecoder(r.Body).Decode(&createUser)
	if err != nil {
		uc.log.Error("decoding error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if createUser.Password != createUser.ConfirmPassword {
		uc.log.Error("passwords do not match")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	users, err := uc.handler.GetList(offset, limit)
	for _, u := range users {
		if u.GetLogin() == createUser.Login {
			uc.log.Error("login already exists")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	hashedPassword, err := hashAndSalt([]byte(createUser.Password))
	if err != nil {
		uc.log.Error("password hashing error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser := user.NewUser(createUser.Login, hashedPassword, createUser.Email)

	uc.handler.Storage.Add(*newUser)

	response.Render(w, uc.log, http.StatusCreated, newUser)
}

func (uc *Controller) GetUserList(w http.ResponseWriter, r *http.Request) {

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	users, err := uc.handler.GetList(limit, offset)

	if err != nil {
		uc.log.Error("getting list error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Render(w, uc.log, http.StatusOK, users)
}

func (uc *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUser user_handler.CreateUserCommand

	err := json.NewDecoder(r.Body).Decode(&createUser)
	if err != nil {
		uc.log.Error("decoding error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uuid := uc.handler.Create(createUser)

	response.Render(w, uc.log, http.StatusCreated, uuid)
}

func (uc *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	err := uc.handler.Delete(userId)
	if err != nil {
		uc.log.Error("deleting error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Render(w, uc.log, http.StatusNoContent, nil)
}

func (uc *Controller) Authenticate(w http.ResponseWriter, r *http.Request) {
	var credentials user_handler.AuthenticateCommand

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		uc.log.Error("decoding error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := uc.handler.Authenticate(credentials.Login, credentials.Password)
	if err != nil {
		uc.log.Error("auth error", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token.JwtToken,
	})

	response.Render(w, uc.log, http.StatusOK, token)
}

func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
