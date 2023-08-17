package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"user-service/internal/controller/user"
	"user-service/internal/persistance/psql"
	"user-service/internal/user/user_handler"
	"user-service/pkg/database"
)

func main() {
	logger, err := zap.NewProduction()
	logger.Info("Starting server...")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cannot init logger: %s", err)
		os.Exit(2)
	}

	config, err := ReadConfig()
	if err != nil {
		logger.Fatal("Cannot read config", zap.Error(err))
	}

	db, err := database.NewConnection(config.PsqlDSN)
	if err != nil {
		logger.Fatal("Cannot init database", zap.Error(err))
	}

	router := chi.NewRouter()

	userStorage := psql.NewUserPsqlStorage(db, logger)
	userHandler := user_handler.NewHandler(userStorage)

	userController := user.NewUserController(logger, userHandler)
	userController.RegisterRoutes(router)

	if err = http.ListenAndServe(":8085", router); err != nil {
		logger.Fatal("Server can't start", zap.Error(err))
		return
	}

}
