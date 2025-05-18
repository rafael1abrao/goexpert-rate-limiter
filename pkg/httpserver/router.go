package httpserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafael1abrao/goexpert-rate-limiter/internal/limiter"
	"github.com/rafael1abrao/goexpert-rate-limiter/internal/middleware"
)

type ServerConfig struct {
	Port    string
	Limiter *limiter.LimiterService
}

func NewServer(cfg ServerConfig) *fiber.App {
	app := fiber.New()

	// Middlewares globais
	app.Use(middleware.NewRateLimiterMiddleware(cfg.Limiter))

	// Rotas
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Rate Limiter is running")
	})

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	return app
}
