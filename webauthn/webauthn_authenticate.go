package webauthn

import (
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/sourcecode081017/passkey-auth-go/internal/models"
)

func WebAuthnAuthBegin(user *models.User) *protocol.CredentialAssertion {
	_webauthn := InitWebAuthn()
	credentials, err := GetUserCredentials(user)
	if err != nil {
		fmt.Printf("Error getting existing credentials: %v", err)
	}
	if len(credentials) == 0 {
		panic("No credentials found for user")
	}
	var allowList []protocol.CredentialDescriptor
	for _, cred := range credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		allowList = append(allowList, descriptor)
	}
	options, sessionData, err := _webauthn.BeginLogin(user, webauthn.WithAllowedCredentials(allowList))
	if err != nil {
		panic(err)
	}
	// save session data to redis
	err = saveSessionData(user, sessionData, "WEBAUTHN_AUTH")
	if err != nil {
		panic(err)
	}

	return options
}

func WebAuthnAuthComplete(r *http.Request, user *models.User) error {
	_webauthn := InitWebAuthn()

	// Retrieve session data from Redis
	sessionData, err := getSessionData(user)
	if err != nil {
		return err
	}
	fmt.Printf("session data: %+v\n", sessionData)

	// Complete the authentication process
	_, err = _webauthn.FinishLogin(user, *sessionData, r)
	if err != nil {
		fmt.Printf("Error finishing authentication: %v\n", err)
		return err
	}

	return nil
}
