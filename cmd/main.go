package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"user-service/internal/controller/user"
	"user-service/internal/controller/user/profile"
	user_storage "user-service/internal/persistance/psql/user"
	profile_storage "user-service/internal/persistance/psql/user/profile"
	"user-service/internal/user/profile/profile_handler"
	"user-service/internal/user/token"
	"user-service/internal/user/user_handler"
	"user-service/pkg/cache"
	"user-service/pkg/database"
)

func main() {
	logger, err := zap.NewProduction()
	logger.Info("Starting server...")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cannot init logger: %s", err)
		os.Exit(2)
	}

	config, err := readConfig()
	if err != nil {
		logger.Fatal("Cannot read config", zap.Error(err))
	}

	db, err := database.NewConnection(config.PsqlDSN)
	if err != nil {
		logger.Fatal("Cannot init database", zap.Error(err))
	}

	cacheProvider, err := cache.NewCache(config.RedisAddr)
	if err != nil {
		logger.Fatal("Cannot init cache", zap.Error(err))
	}

	jwtGenerator := token.NewJwtGenerator(config.JwtSecret)
	cachedJwtGenerator := token.NewGeneratorCached(jwtGenerator, cacheProvider)

	router := chi.NewRouter()

	userStorage := user_storage.NewUserPsqlStorage(db, logger)
	userHandler := user_handler.NewHandler(userStorage, cachedJwtGenerator)

	userController := user.NewUserController(userHandler, logger)
	userController.RegisterRoutes(router)

	profileStorage := profile_storage.NewProfilePsqlStorage(db, logger)
	profileHandler := profile_handler.NewProfileHandler(profileStorage)

	profileController := profile.NewProfileController(profileHandler, logger)
	profileController.RegisterRoutes(router)

	if err = http.ListenAndServe(":8085", router); err != nil {
		logger.Fatal("Server can't start", zap.Error(err))
		return
	}

}
