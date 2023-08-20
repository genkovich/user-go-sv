package profile

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"user-service/internal/user/profile/profile_handler"
	"user-service/pkg/response"
)

type Controller struct {
	handler *profile_handler.Handler
	log     *zap.Logger
}

func (uc *Controller) RegisterRoutes(r *chi.Mux) {
	r.Get("/users/{userId}/profile", uc.GetUserProfile)
	r.Post("/users/{userId}/profile", uc.UpdateUserProfile)
}

func NewProfileController(handler *profile_handler.Handler, log *zap.Logger) *Controller {
	return &Controller{
		handler: handler,
		log:     log,
	}
}

func (uc *Controller) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	userUuid, err := uuid.Parse(userId)

	if err != nil {
		uc.log.Error("failed parse uuid", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profile, err := uc.handler.GetProfile(userUuid)

	if err != nil {
		uc.log.Error("failed get profile", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Render(w, uc.log, http.StatusOK, profile)
}

func (uc *Controller) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	userUuid, err := uuid.Parse(userId)

	if err != nil {
		uc.log.Error("failed parse uuid", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var updateProfile profile_handler.UpdateProfileCommand

	err = json.NewDecoder(r.Body).Decode(&updateProfile)
	if err != nil {
		uc.log.Error("decoding error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updateProfile.UserId = userUuid
	validate := validator.New()
	err = validate.Struct(updateProfile)
	if err != nil {
		var errors []string
		for _, tmpError := range err.(validator.ValidationErrors) {
			errors = append(errors, tmpError.Error())
		}
		response.Render(w, uc.log, http.StatusBadRequest, errors)
		return
	}

	err = uc.handler.UpdateProfile(updateProfile)

	if err != nil {
		uc.log.Error("failed update profile", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Render(w, uc.log, http.StatusOK, nil)
}
