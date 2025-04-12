package webauthn

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sourcecode081017/passkey-auth-go/internal/cache"
)

func CheckUserExists(ctx context.Context, username string) (bool, error) {
	// check if user keys exist in redis
	redisCache := ctx.Value("cache").(*cache.RedisCache)
	pattern := fmt.Sprintf("WEBAUTHN_CREDENTIAL_%s*", username)
	users, err := redisCache.Keys(ctx, pattern)
	if err != nil {
		return false, fmt.Errorf("failed to get keys from Redis: %v", err)
	}
	if len(users) == 0 {
		return false, nil
	}
	return true, nil
}

func DeleteUserKey(ctx context.Context, username string, credentialId string) error {
	// delete user keys in redis
	redisCache := ctx.Value("cache").(*cache.RedisCache)
	key := fmt.Sprintf("WEBAUTHN_CREDENTIAL_%s_%s", username, credentialId)
	err := redisCache.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete key from Redis: %v", err)
	}
	return nil
}

func UpdateLastUsedAt(ctx context.Context, username string, credentialId string) error {
	// update last used at in redis
	redisCache := ctx.Value("cache").(*cache.RedisCache)
	key := fmt.Sprintf("WEBAUTHN_CREDENTIAL_%s_%s", username, credentialId)
	credentialJSON, err := redisCache.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to get key from Redis: %v", err)
	}
	var passkeyCredential PasskeyCredential
	err = json.Unmarshal([]byte(credentialJSON), &passkeyCredential)
	if err != nil {
		return fmt.Errorf("failed to unmarshal credential data: %v", err)
	}
	newPasskeyCredential := &PasskeyCredential{
		Credential: passkeyCredential.Credential,
		CreatedAt:  passkeyCredential.CreatedAt,
		LastUsedAt: func() *time.Time {
			t := time.Now()
			return &t
		}(),
	}
	credentialJSON, err = json.Marshal(*newPasskeyCredential)
	if err != nil {
		return fmt.Errorf("failed to marshal credential data: %v", err)
	}
	err = redisCache.Set(ctx, key, credentialJSON, 0)
	if err != nil {
		return fmt.Errorf("failed to set value in Redis: %v", err)
	}
	return nil
}
