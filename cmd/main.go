package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sourcecode081017/passkey-auth-go/internal/cache"
	"github.com/sourcecode081017/passkey-auth-go/internal/rest"
)

func main() {
	fmt.Println("Starting HTTP server on port 8080")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" || redisPort == "" {
		log.Fatalf("REDIS_HOST or REDIS_PORT environment variables are not set")
	}
	redisURL := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisCache, err := connectToRedisWithRetry(redisURL, "", 0)
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		os.Exit(1)
	}
	// After creating redisCache
	if err := redisCache.Ping(); err != nil {
		fmt.Printf("Redis connection test failed: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Redis connection test successful")
	}
	restHandler := rest.NewRestHandler(redisCache)
	restHandler.StartHttpServer()
}

// Add this function to your code
func connectToRedisWithRetry(redisURL string, password string, db int) (*cache.RedisCache, error) {
	var redisCache *cache.RedisCache
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		redisCache, err = cache.NewRedisCache(redisURL, password, db)
		if err == nil {
			fmt.Println("Successfully connected to Redis")
			return redisCache, nil
		}

		fmt.Printf("Failed to connect to Redis (attempt %d/%d): %v\n", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(3 * time.Second)
		}
	}
	return nil, err
}
