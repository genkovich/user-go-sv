package main

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"user-service/internal/controller/user"
)

func main() {
	router := chi.NewRouter()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting server...")

	userController := user.NewUserController(logger)
	userController.RegisterRoutes(router)

	if err := http.ListenAndServe(":8085", router); err != nil {
		logger.Fatal("Server can't start", zap.Error(err))
		return
	}
}

