package webauthn

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/sourcecode081017/passkey-auth-go/internal/cache"
	"github.com/sourcecode081017/passkey-auth-go/internal/models"
)

func WebAuthnRegisterBegin(user *models.User) *protocol.CredentialCreation {
	webauthn := InitWebAuthn()
	options, sessionData, err := webauthn.BeginRegistration(user)
	if err != nil {
		panic(err)
	}
	fmt.Println("SessionData: ", sessionData)
	// save session data to redis
	err = saveSessionData(user, sessionData)
	if err != nil {
		panic(err)
	}

	return options

}

func WebAuthnRegisterComplete(r *http.Request, user *models.User) error {
	webauthn := InitWebAuthn()

	// Retrieve session data from Redis
	sessionData, err := getSessionData(user)
	if err != nil {
		return err
	}
	fmt.Printf("session data: %+v\n", sessionData)

	// Complete the registration process
	credential, err := webauthn.FinishRegistration(user, *sessionData, r)
	if err != nil {
		fmt.Printf("Error finishing registration: %v\n", err)
		return err
	}
	// print the credential for debugging
	fmt.Printf("Registered Credential: %+v\n", credential)

	// Registration successful, you can now store the credential in your database

	return nil
}

func saveSessionData(user *models.User, sessionData *webauthn.SessionData) error {
	redis_cache, err := cache.NewRedisCache("localhost:6379", "", 0)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return err
	}
	defer redis_cache.Close()

	ctx := context.Background()

	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		log.Fatalf("Failed to marshal session data: %v", err)
		return err
	}
	redisKey := fmt.Sprintf("WEBAUTHN_REGISTER_%s", user.Username)

	// Set a value
	err = redis_cache.Set(ctx, redisKey, sessionDataJSON, 60*time.Second)
	if err != nil {
		log.Fatalf("Failed to set value in Redis: %v", err)
		return err
	}
	return nil
}

func getSessionData(user *models.User) (*webauthn.SessionData, error) {
	redis_cache, err := cache.NewRedisCache("localhost:6379", "", 0)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil, err
	}
	defer redis_cache.Close()

	ctx := context.Background()
	redisKey := fmt.Sprintf("WEBAUTHN_REGISTER_%s", user.Username)

	// Get the value from Redis
	sessionDataJSON, err := redis_cache.Get(ctx, redisKey)
	if err != nil {
		log.Fatalf("Failed to get value from Redis: %v", err)
		return nil, err
	}
	if sessionDataJSON == "" {
		return nil, fmt.Errorf("session data not found")
	}

	var sessionData webauthn.SessionData
	err = json.Unmarshal([]byte(sessionDataJSON), &sessionData)
	if err != nil {
		log.Fatalf("Failed to unmarshal session data: %v", err)
		return nil, err
	}
	return &sessionData, nil
}
