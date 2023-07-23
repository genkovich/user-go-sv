package user

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"user-service/internal/persistance/memory"
	"user-service/internal/user/user_handler"
	"user-service/pkg/response"
)

type Controller struct {
	handler user_handler.Handler
	log     *zap.Logger
}

func NewUserController(log *zap.Logger) *Controller {
	return &Controller{
		handler: user_handler.Handler{
			Storage: memory.NewUserMemoryStorage(),
		},
		log: log,
	}
}

func (uc *Controller) RegisterRoutes(r *chi.Mux) {
	r.Get("/user", uc.GetUserList)
	r.Post("/user", uc.CreateUser)
	r.Delete("/user/{userId}", uc.DeleteUser)
}

func (uc *Controller) GetUserList(w http.ResponseWriter, r *http.Request) {
	users := uc.handler.GetList()

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
