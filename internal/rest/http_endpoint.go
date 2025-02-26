package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sourcecode081017/passkey-auth-go/internal/models"
	"github.com/sourcecode081017/passkey-auth-go/webauthn"
)

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func registerInitiate(c *gin.Context) {
	username := c.Param("username")
	webauthnUser := *models.GetUser(username)
	optionsData := webauthn.WebAuthnRegisterBegin(&webauthnUser)
	fmt.Println("Options: ", optionsData)
	c.JSON(200, gin.H{
		"message": "OK",
		"options": optionsData,
	})

}

func StartHttpServer() {

	router := gin.Default()
	// Global endpoint
	router.GET("/", healthCheck)
	router.POST("passkey-auth/register-initiate/:username", registerInitiate)
	// Start the HTTP server with GIN
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
