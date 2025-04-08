package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"restaurant-management/internal/config"
	"restaurant-management/internal/delivery/http"
	"restaurant-management/internal/repository"
	"restaurant-management/internal/repository/postgres"
	"restaurant-management/internal/usecase"
	"restaurant-management/pkg/database"
	// Раскомментируйте для генерации документации Swagger
	// _ "restaurant-management/docs"
)

// @title Restaurant Management System API
// @version 1.0
// @description API сервера для системы управления рестораном
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {
	configPath := getConfigPath()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}
	defer db.Close()

	repos := initRepositories(db)
	useCases := initUseCases(repos)

	server := http.NewServer(cfg, useCases)
	log.Printf("Сервер запущен на порту %s", cfg.Server.Port)
	if err := server.Start(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = filepath.Join("configs", "config.yaml")
	}
	return configPath
}

func initDB(cfg *config.Config) (*database.PostgreSQL, error) {
	ctx := context.Background()
	db, err := database.NewPostgreSQL(ctx, cfg.Database.PostgresURL())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initRepositories(db *database.PostgreSQL) *repository.Repository {
	return &repository.Repository{
		User:       postgres.NewUserRepository(db.Pool),
		City:       postgres.NewCityRepository(db.Pool),
		Restaurant: postgres.NewRestaurantRepository(db.Pool),
	}
}

func initUseCases(repos *repository.Repository) *usecase.UseCase {
	return &usecase.UseCase{
		User:       usecase.NewUserUseCase(repos.User),
		City:       usecase.NewCityUseCase(repos.City),
		Restaurant: usecase.NewRestaurantUseCase(repos.Restaurant, repos.City),
	}
}
