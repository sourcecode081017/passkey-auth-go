package webauthn

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
