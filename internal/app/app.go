package app

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
)

type App struct {
	config  *config.Config
	db      *database.PostgreSQL
	server  *http.Server
	useCase *usecase.UseCase
	repos   *repository.Repository
}

func New() (*App, error) {
	app := &App{}

	configPath := getConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	app.config = cfg

	db, err := initDB(cfg)
	if err != nil {
		return nil, err
	}
	app.db = db

	app.repos = initRepositories(db)

	app.useCase = initUseCases(app.repos)

	app.server = http.NewServer(cfg, app.useCase)

	return app, nil
}

func (a *App) Run() error {
	log.Printf("Сервер запущен на порту %s", a.config.Server.Port)
	return a.server.Start()
}

func (a *App) Stop() {
	if a.db != nil {
		a.db.Close()
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
