package webauthn

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/sourcecode081017/passkey-auth-go/internal/cache"
	"github.com/sourcecode081017/passkey-auth-go/internal/models"
)

type PasskeyCredential struct {
	Credential *webauthn.Credential
	CreatedAt  *time.Time `json:"createdAt"`
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
}

func GetWebAuthnKeys(ctx context.Context, username string) ([]PasskeyCredential, error) {
	// Get the user from the database
	user := models.GetUser(username)
	passkeys, err := GetPasskeysForUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error getting user credentials: %v", err)
	}
	return passkeys, nil
}

func GetPasskeysForUser(ctx context.Context, user *models.User) ([]PasskeyCredential, error) {
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

	var passkeys []PasskeyCredential

	// Retrieve and unmarshal each credential
	for _, key := range keys {
		credentialJSON, err := redisCache.Get(ctx, key)
		if err != nil {
			log.Printf("Failed to get credential for key %s: %v", key, err)
			continue
		}

		var passkeyCredential PasskeyCredential
		if err := json.Unmarshal([]byte(credentialJSON), &passkeyCredential); err != nil {
			log.Printf("Failed to unmarshal credential data: %v", err)
			continue
		}

		passkeys = append(passkeys, passkeyCredential)
	}

	return passkeys, nil

}
