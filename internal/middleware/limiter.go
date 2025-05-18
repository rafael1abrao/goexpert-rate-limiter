package middleware

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rafael1abrao/goexpert-rate-limiter/internal/limiter"
)

func NewRateLimiterMiddleware(s *limiter.LimiterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()

		ip := c.IP()
		token := c.Get("API_KEY")

		allowed, err := s.IsRequestAllowed(ctx, ip, token)
		if err != nil {
			log.Printf("Rate limiter error: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
		}

		if !allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
		}

		return c.Next()
	}
}
