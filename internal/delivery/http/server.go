package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"restaurant-management/internal/config"
	"restaurant-management/internal/delivery/http/handlers"
	"restaurant-management/internal/usecase"
)

type Server struct {
	echo    *echo.Echo
	config  *config.Config
	useCase *usecase.UseCase
}

func NewServer(cfg *config.Config, useCase *usecase.UseCase) *Server {
	e := echo.New()
	return &Server{
		echo:    e,
		config:  cfg,
		useCase: useCase,
	}
}

func (s *Server) Start() error {
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.CORS())

	s.setupRoutes()

	go func() {
		address := fmt.Sprintf(":%s", s.config.Server.Port)
		if err := s.echo.Start(address); err != nil && err != http.ErrServerClosed {
			s.echo.Logger.Fatalf("Ошибка при старте сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.echo.Shutdown(ctx)
}

func (s *Server) setupRoutes() {
	s.echo.GET("/swagger*", echoSwagger.WrapHandler)

	api := s.echo.Group("/api/v1")

	userHandler := handlers.NewUserHandler(s.useCase.User)
	userHandler.Register(api)

	cityHandler := handlers.NewCityHandler(s.useCase.City)
	cityHandler.Register(api)

	restaurantHandler := handlers.NewRestaurantHandler(s.useCase.Restaurant)
	restaurantHandler.Register(api)

	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}
