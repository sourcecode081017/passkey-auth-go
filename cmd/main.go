package main

import (
	"fmt"
	"os"

	"github.com/sourcecode081017/passkey-auth-go/internal/cache"
	"github.com/sourcecode081017/passkey-auth-go/internal/rest"
)

func main() {
	fmt.Println("Starting HTTP server on port 8080")
	redisCache, err := cache.NewRedisCache("localhost:6379", "", 0)
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		os.Exit(1)
	}
	restHandler := rest.NewRestHandler(redisCache)
	restHandler.StartHttpServer()
}
