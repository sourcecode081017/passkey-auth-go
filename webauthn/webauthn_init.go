package webauthn

import "github.com/go-webauthn/webauthn/webauthn"

func InitWebAuthn() *webauthn.WebAuthn {
	wautn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Passkey Auth Go",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:8080"},
	})
	if err != nil {
		panic(err)
	}
	return wautn
}
