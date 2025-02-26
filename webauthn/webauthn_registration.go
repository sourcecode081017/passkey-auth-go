package webauthn

import (
	"fmt"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/sourcecode081017/passkey-auth-go/internal/models"
)

func WebAuthnRegisterBegin(user *models.User) *protocol.CredentialCreation {
	webauthn := InitWebAuthn()
	options, sessionData, err := webauthn.BeginRegistration(user)
	if err != nil {
		panic(err)
	}
	fmt.Println("SessionData: ", sessionData)
	return options

}
