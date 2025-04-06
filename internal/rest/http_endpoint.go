package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sourcecode081017/passkey-auth-go/internal/cache"
	"github.com/sourcecode081017/passkey-auth-go/internal/middleware"
	"github.com/sourcecode081017/passkey-auth-go/internal/models"
	"github.com/sourcecode081017/passkey-auth-go/webauthn"
)

type RestHandler struct {
	redisCache *cache.RedisCache
}

func NewRestHandler(redisCache *cache.RedisCache) *RestHandler {
	return &RestHandler{
		redisCache: redisCache,
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func (r *RestHandler) registerInitiate(c *gin.Context) {
	username := c.Param("username")
	webauthnUser := *models.GetUser(username)
	redisCache := r.redisCache
	// set cache object in context
	c.Set("cache", redisCache)
	optionsData := webauthn.WebAuthnRegisterBegin(c, &webauthnUser)
	fmt.Println("Options: ", optionsData)
	c.JSON(200, gin.H{
		"message": "OK",
		"options": optionsData,
	})

}

func (r *RestHandler) registerComplete(c *gin.Context) {
	username := c.Param("username")
	webauthnUser := *models.GetUser(username)
	// Pass the HTTP request to WebAuthnRegisterComplete
	httpRequest := c.Request
	redisCache := r.redisCache
	// set cache object in context
	c.Set("cache", redisCache)
	err := webauthn.WebAuthnRegisterComplete(c, httpRequest, &webauthnUser)
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

func (r *RestHandler) authInitiate(c *gin.Context) {
	username := c.Param("username")
	redisCache := r.redisCache
	// set cache object in context
	c.Set("cache", redisCache)
	exists, err := webauthn.CheckUserExists(c, username)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !exists {
		c.JSON(404, gin.H{
			"error": "User not found",
		})
		return
	}
	webauthnUser := *models.GetUser(username)
	optionsData := webauthn.WebAuthnAuthBegin(c, &webauthnUser)
	fmt.Println("Options: ", optionsData)
	c.JSON(200, gin.H{
		"message": "OK",
		"options": optionsData,
	})
}

func (r *RestHandler) authComplete(c *gin.Context) {
	username := c.Param("username")
	webauthnUser := *models.GetUser(username)
	// Pass the HTTP request to WebAuthnAuthComplete
	httpRequest := c.Request
	redisCache := r.redisCache
	// set cache object in context
	c.Set("cache", redisCache)
	err := webauthn.WebAuthnAuthComplete(c, httpRequest, &webauthnUser)
	// Verify the authentication response
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Authentication successful",
	})
}

func (r *RestHandler) getPassKeys(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(400, gin.H{
			"error": "Username is required",
		})
		return
	}
	redisCache := r.redisCache
	// set cache object in context
	c.Set("cache", redisCache)
	keys, err := webauthn.GetWebAuthnKeys(c, username)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to retrieve keys",
		})
		return
	}

	c.JSON(200, gin.H{
		"keys": keys,
	})
}

func (r *RestHandler) deleteUserKey(c *gin.Context) {
	username := c.Param("username")
	// get credentialId from body
	var requestBody struct {
		CredentialId string `json:"passkeyId"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	credentialId := requestBody.CredentialId
	if username == "" || credentialId == "" {
		c.JSON(400, gin.H{
			"error": "Username and credentialId are required",
		})
		return
	}
	redisCache := r.redisCache
	// set cache object in context
	c.Set("cache", redisCache)
	err := webauthn.DeleteUserKey(c, username, credentialId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to delete key",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Key deleted successfully",
	})
}

func (r *RestHandler) StartHttpServer() {

	router := gin.Default()
	// Global endpoint
	router.Use(middleware.CORSMiddleware())
	router.GET("/", healthCheck)
	router.POST("passkey-auth/register-initiate/:username", r.registerInitiate)
	router.POST("passkey-auth/register-complete/:username", r.registerComplete)
	router.POST("passkey-auth/auth-initiate/:username", r.authInitiate)
	router.POST("passkey-auth/auth-complete/:username", r.authComplete)
	router.GET("passkey-auth/:username/keys", r.getPassKeys)
	router.DELETE("passkey-auth/:username/keys", r.deleteUserKey)
	// Start the HTTP server with GIN
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
