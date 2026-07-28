package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"job4j.ru/sharetrip-notification/configs"
	"job4j.ru/sharetrip-notification/internal/api"
	"job4j.ru/sharetrip-notification/internal/domain"
	"job4j.ru/sharetrip-notification/internal/repository"
	"job4j.ru/sharetrip-notification/internal/service"
	"log"
)

func main() {
	ctx := context.Background()

	cfg := repository.Config{
		Host:     configs.Env("DB_HOST", "localhost"),
		Port:     configs.EnvInt("DB_PORT", 6544),
		User:     configs.Env("DB_USER", "postgres"),
		Password: configs.Env("DB_PASSWORD", "password"),
		DBName:   configs.Env("DB_NAME", "notifications"),
		SSLMode:  configs.Env("DB_SSLMODE", "disable"),
	}

	pool, err := repository.NewPool(ctx, cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = pool.Close() }()

	repo := repository.NewRepoPg(pool)

	server := api.NewServer(repo)
	server.NotificationService = &service.NotificationService{
		Pool: pool,
		NotificationUsecase: &domain.NotificationUsecase{
			NotificationRepo: repo,
		},
	}

	app := fiber.New()
	server.Route(app.Group("/"))

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
