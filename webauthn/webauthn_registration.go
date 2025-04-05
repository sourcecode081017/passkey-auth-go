package webauthn

import (
	"context"
	"encoding/base64"
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

func WebAuthnRegisterBegin(ctx context.Context, user *models.User) *protocol.CredentialCreation {
	_webauthn := InitWebAuthn()
	existingCredentials, err := GetUserCredentials(ctx, user)
	if err != nil {
		fmt.Printf("Error getting existing credentials: %v", err)
	}
	var excludeList []protocol.CredentialDescriptor
	for _, cred := range existingCredentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		excludeList = append(excludeList, descriptor)
	}

	options, sessionData, err := _webauthn.BeginRegistration(user, webauthn.WithExclusions(excludeList))
	if err != nil {
		panic(err)
	}
	// save session data to redis
	err = saveSessionData(ctx, user, sessionData, "WEBAUTHN_REGISTER")
	if err != nil {
		panic(err)
	}

	return options

}

func WebAuthnRegisterComplete(ctx context.Context, r *http.Request, user *models.User) error {
	webauthn := InitWebAuthn()

	// Retrieve session data from Redis
	sessionData, err := getSessionData(ctx, user, "WEBAUTHN_REGISTER")
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
	// save credential data to redis
	err = saveCredentialData(ctx, user, credential)
	if err != nil {
		fmt.Printf("Error saving credential data: %v\n", err)
		return err
	}

	return nil
}

func saveSessionData(ctx context.Context, user *models.User, sessionData *webauthn.SessionData, prefix string) error {
	redis_cache, err := cache.NewRedisCache("localhost:6379", "", 0)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return err
	}
	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		log.Fatalf("Failed to marshal session data: %v", err)
		return err
	}
	redisKey := fmt.Sprintf("%s_%s", prefix, user.Username)

	// Set a value
	err = redis_cache.Set(ctx, redisKey, sessionDataJSON, 60*time.Second)
	if err != nil {
		log.Fatalf("Failed to set value in Redis: %v", err)
		return err
	}
	return nil
}

func getSessionData(ctx context.Context, user *models.User, prefix string) (*webauthn.SessionData, error) {
	redisCache := ctx.Value("cache").(*cache.RedisCache)
	if redisCache == nil {
		return nil, fmt.Errorf("cache not found in context")
	}
	redisKey := fmt.Sprintf("%s_%s", prefix, user.Username)

	// Get the value from Redis
	sessionDataJSON, err := redisCache.Get(ctx, redisKey)
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

func saveCredentialData(ctx context.Context, user *models.User, credential *webauthn.Credential) error {
	redisCache := ctx.Value("cache").(*cache.RedisCache)
	if redisCache == nil {
		return fmt.Errorf("cache not found in context")
	}

	// Base64 encode the credential ID for use in the Redis key
	credentialIDBase64 := base64.StdEncoding.EncodeToString(credential.ID)

	credentialJSON, err := json.Marshal(credential)
	if err != nil {
		log.Printf("Failed to marshal credential data: %v", err)
		return err
	}

	// Use the base64 encoded credential ID in the Redis key
	redisKey := fmt.Sprintf("WEBAUTHN_CREDENTIAL_%s_%s", user.Username, credentialIDBase64)

	// Set a value
	err = redisCache.Set(ctx, redisKey, credentialJSON, 0)
	if err != nil {
		log.Printf("Failed to set value in Redis: %v", err)
		return err
	}

	return nil
}

func GetUserCredentials(ctx context.Context, user *models.User) ([]webauthn.Credential, error) {
	redisCache := ctx.Value("cache").(*cache.RedisCache)
	if redisCache == nil {
		return nil, fmt.Errorf("cache not found in context")
	}
	// Pattern to match all credentials for this user
	pattern := fmt.Sprintf("WEBAUTHN_CREDENTIAL_%s_*", user.Username)

	// Get all keys matching the pattern
	keys, err := redisCache.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential keys from Redis: %v", err)
	}

	var credentials []webauthn.Credential

	// Retrieve and unmarshal each credential
	for _, key := range keys {
		credentialJSON, err := redisCache.Get(ctx, key)
		if err != nil {
			log.Printf("Failed to get credential for key %s: %v", key, err)
			continue
		}

		var credential webauthn.Credential
		if err := json.Unmarshal([]byte(credentialJSON), &credential); err != nil {
			log.Printf("Failed to unmarshal credential data: %v", err)
			continue
		}

		credentials = append(credentials, credential)
	}

	return credentials, nil
}
