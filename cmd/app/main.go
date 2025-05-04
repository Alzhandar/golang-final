package main

import (
	"log"

	_ "restaurant-management/docs"
	"restaurant-management/internal/app"
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
	application, err := app.New()
	if err != nil {
		log.Fatalf("Ошибка при инициализации приложения: %v", err)
	}
	defer application.Stop()

	if err := application.Run(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
