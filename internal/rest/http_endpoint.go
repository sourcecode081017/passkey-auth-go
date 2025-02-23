package rest

import (
	"github.com/gin-gonic/gin"
)

func StartHttpServer() {
	// Start the HTTP server with GIN
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
