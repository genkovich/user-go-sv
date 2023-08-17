package user

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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

func NewUserController(log *zap.Logger, handler *user_handler.Handler) *Controller {
	return &Controller{
		handler: handler,
		log:     log,
	}
}

func (uc *Controller) RegisterRoutes(r *chi.Mux) {
	r.Get("/user", uc.GetUserList)
	r.Post("/user", uc.CreateUser)
	r.Delete("/user/{userId}", uc.DeleteUser)
	r.Post("/user/register", uc.RegisterUser)
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

	users, _ := uc.handler.GetList(offset, limit)
	for _, u := range users {
		if u.GetLogin() == createUser.Login {
			uc.log.Error("login already exists")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	newUser := user.CreateUser(createUser.Login, createUser.Password)
	newUser.Email = createUser.Email

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
