package main

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"github.com/rafael1abrao/goexpert-rate-limiter/internal/limiter"
	redislimiter "github.com/rafael1abrao/goexpert-rate-limiter/internal/limiter/redis"
	"github.com/rafael1abrao/goexpert-rate-limiter/pkg/httpserver"
)

func main() {
	_ = godotenv.Load()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	strategy := redislimiter.NewRedisStrategy(redisClient)
	service := limiter.NewLimiterService(strategy)

	app := httpserver.NewServer(httpserver.ServerConfig{
		Port:    "8080",
		Limiter: service,
	})

	log.Fatal(app.Listen(":8080"))
}
