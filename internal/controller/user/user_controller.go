package user

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
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
	r.Post("/users/change-password", uc.ChangeUserPassword)
	r.Post("/change-role", uc.ChangeUserRole)
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
		limit = 0
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

	_, err = uc.handler.Storage.Add(*newUser)
	if err != nil {
		uc.log.Error("Error adding new user to the database", zap.Error(err))
		response.Render(w, uc.log, http.StatusInternalServerError, nil)
		return
	}

	uc.log.Info("New user added to the database", zap.String("Login", newUser.GetLogin()))
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

func (uc *Controller) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	var changePasswordCommand struct {
		UserId      string `json:"user_id"`
		NewPassword string `json:"new_password"`
	}

	err := json.NewDecoder(r.Body).Decode(&changePasswordCommand)
	if err != nil {
		response.Render(w, uc.log, http.StatusBadRequest, nil)
		return
	}

	targetUser, err := uc.handler.Storage.GetByLogin(changePasswordCommand.UserId)
	if err != nil {
		response.Render(w, uc.log, http.StatusNotFound, nil)
		return
	}

	err = uc.handler.Storage.UpdatePassword(targetUser.GetId().String(), changePasswordCommand.NewPassword)
	if err != nil {
		response.Render(w, uc.log, http.StatusInternalServerError, nil)
		return
	}

	response.Render(w, uc.log, http.StatusOK, nil)
}

func (uc *Controller) ChangeUserRole(w http.ResponseWriter, r *http.Request) {

	var changeRoleCommand struct {
		UserId  string `json:"user_id"`
		NewRole string `json:"new_role"`
	}

	err := json.NewDecoder(r.Body).Decode(&changeRoleCommand)
	if err != nil {
		uc.log.Error("decoding error", zap.Error(err))
		response.Render(w, uc.log, http.StatusBadRequest, nil)
		return
	}

	claims, err := parseJWT(r)
	if err != nil {
		uc.log.Error("authentication error", zap.Error(err))
		response.Render(w, uc.log, http.StatusUnauthorized, nil)
		return
	}

	userRole, ok := claims["role"].(string)
	if !ok || userRole != "ROLE_ADMIN" {
		uc.log.Warn("unauthorized access", zap.String("user_role", userRole))
		response.Render(w, uc.log, http.StatusForbidden, nil)
		return
	}

	targetUser, err := uc.handler.Storage.GetByLogin(changeRoleCommand.UserId)
	if err != nil {
		response.Render(w, uc.log, http.StatusNotFound, nil)
		return
	}

	err = uc.handler.Storage.UpdateUserRole(targetUser.GetId().String(), changeRoleCommand.NewRole)
	if err != nil {
		response.Render(w, uc.log, http.StatusInternalServerError, nil)
		return
	}

	response.Render(w, uc.log, http.StatusOK, nil)
}

func parseJWT(r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, errors.New("token not found")
		}
		return nil, err
	}

	tokenString := cookie.Value

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
