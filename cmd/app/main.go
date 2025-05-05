package main

import (
	"context"
	"log"
	"time"

	_ "restaurant-management/docs"
	"restaurant-management/internal/app"
	"restaurant-management/internal/config"
	"restaurant-management/internal/iiko"

	"github.com/go-redis/redis/v8"
)

func main() {
	iikoConfig := config.NewDefaultConfig()

	redisOptions := &redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	}

	redisClient := redis.NewClient(redisOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Предупреждение: Не удалось подключиться к Redis: %v", err)
		log.Printf("Приложение будет запущено без поддержки Redis и IIKO")
	} else {
		iikoService := iiko.NewIikoService(iikoConfig.APILogin, redisClient)
		waiterService := iiko.NewIikoWaiterService(iikoConfig.WaiterAPIURL, iikoConfig.WaiterAPIKey)
		log.Printf("Инициализированы сервисы IIKO: %v, %v", iikoService != nil, waiterService != nil)
		defer redisClient.Close()
	}

	application, err := app.New()
	if err != nil {
		log.Fatalf("Ошибка при инициализации приложения: %v", err)
	}
	defer application.Stop()

	if err := application.Run(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
