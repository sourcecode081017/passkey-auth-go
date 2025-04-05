package webauthn

import (
	"context"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/sourcecode081017/passkey-auth-go/internal/models"
)

func GetWebAuthnKeys(ctx context.Context, username string) ([]webauthn.Credential, error) {
	// Get the user from the database
	user := models.GetUser(username)
	credentials, err := GetUserCredentials(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error getting user credentials: %v", err)
	}
	return credentials, nil
}
