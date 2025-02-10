package main

import (
	logger "github.com/sirupsen/logrus"

	"service-user/internal/app/delivery/http"
	"service-user/internal/app/repository"
	"service-user/internal/app/service"
	"service-user/internal/app/utils"
	"service-user/internal/configs"
	"service-user/internal/server"
	"service-user/pkg/db"

	logging "service-user/pkg/logger"
)

// @title Profile Service
// @version 1.0
// @description Сервис для управления профилем пользователя, включающий создание, обновление, удаление и получение профиля.

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.cookie CookieAuth
// @name access_token
func main() {
	// Загружаем конфигурацию
	cfg, err := configs.LoadConfig("./internal/configs")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	// Настройка логгера
	logging.SetupLogger(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.OutputFile)

	// Подключение к базе данных
	dbConn, err := db.ConnectPostgres(cfg.Database.Dsn)
	if err != nil {
		logger.Fatalf("Database connection failed: %v", err)
	}
	defer dbConn.Close()

	// public key для проверки токена
	auth, err := utils.NewJWTManager(cfg.Auth.PublicKey)
	if err != nil {
		logger.Fatalf("Error creating JWT Manager: %v", err)
		return
	}

	// db.ApplyMigrations(cfg.Database.Dsn, cfg.Database.MigratePath)

	repo := repository.NewRepository(dbConn)
	services := service.NewService(repo)
	handlers := http.NewHandler(services, auth)

	// Настройка и запуск сервера
	server.SetupAndRunServer(&cfg.Server, handlers.InitRoutes())
}
