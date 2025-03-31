package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sourcecode081017/passkey-auth-go/internal/middleware"
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

func registerComplete(c *gin.Context) {
	username := c.Param("username")
	webauthnUser := *models.GetUser(username)
	// Pass the HTTP request to WebAuthnRegisterComplete
	httpRequest := c.Request
	err := webauthn.WebAuthnRegisterComplete(httpRequest, &webauthnUser)
	// Verify the registration response
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Registration successful",
	})
}

func StartHttpServer() {

	router := gin.Default()
	// Global endpoint
	router.Use(middleware.CORSMiddleware())
	router.GET("/", healthCheck)
	router.POST("passkey-auth/register-initiate/:username", registerInitiate)
	router.POST("passkey-auth/register-complete/:username", registerComplete)
	// Start the HTTP server with GIN
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
