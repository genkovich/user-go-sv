package user

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"user-service/internal/user/user_handler"
	"user-service/pkg/response"
)

type Controller struct {
	handler *user_handler.Handler
	log     *zap.Logger
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

	uuid, err := uc.handler.Create(createUser)
	if err != nil {
		uc.log.Error("creating error", zap.Error(err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

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
		Value: token.GetJwtToken(),
	})

	result := struct {
		Success bool   `json:"success"`
		Token   string `json:"token"`
	}{
		Success: true,
		Token:   token.GetJwtToken(),
	}

	response.Render(w, uc.log, http.StatusOK, result)
}
