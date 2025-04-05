package webauthn

import (
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/sourcecode081017/passkey-auth-go/internal/models"
)

func GetWebAuthnKeys(username string) ([]webauthn.Credential, error) {
	// Get the user from the database
	user := models.GetUser(username)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	credentials, err := GetUserCredentials(user)
	if err != nil {
		return nil, fmt.Errorf("error getting user credentials: %v", err)
	}
	return credentials, nil
}
