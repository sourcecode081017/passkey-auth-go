package webauthn

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/sourcecode081017/passkey-auth-go/internal/models"
)

func WebAuthnAuthBegin(ctx context.Context, user *models.User) *protocol.CredentialAssertion {
	_webauthn := InitWebAuthn()
	credentials, err := GetUserCredentials(ctx, user)
	if err != nil {
		fmt.Printf("Error getting existing credentials: %v", err)
	}
	user.Credentials = credentials // ensure the user has the latest credentials
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
	credentialAssertion, sessionData, err := _webauthn.BeginLogin(user, webauthn.WithAllowedCredentials(allowList))
	if err != nil {
		panic(err)
	}
	// save session data to redis
	err = saveSessionData(ctx, user, sessionData, "WEBAUTHN_AUTH")
	if err != nil {
		panic(err)
	}

	return credentialAssertion
}

func WebAuthnAuthComplete(ctx context.Context, r *http.Request, user *models.User) error {
	_webauthn := InitWebAuthn()
	credentials, err := GetUserCredentials(ctx, user)
	if err != nil {
		fmt.Printf("Error getting existing credentials: %v", err)
	}
	user.Credentials = credentials // ensure the user has the latest credentials
	if len(credentials) == 0 {
		panic("No credentials found for user")
	}
	// Retrieve session data from Redis
	sessionData, err := getSessionData(ctx, user, "WEBAUTHN_AUTH")
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
