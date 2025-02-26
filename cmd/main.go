package main

import (
	"fmt"

	"github.com/sourcecode081017/passkey-auth-go/internal/rest"
)

func main() {
	fmt.Println("Starting HTTP server on port 8080")
	//dbPool := db.GetConnectionPool()
	rest.StartHttpServer()
}
