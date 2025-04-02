package models

import (
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	Id          []byte
	Username    string
	DisplayName string
	Credentials []webauthn.Credential
}

func (u *User) WebAuthnID() []byte {
	return u.Id
}

func (u *User) WebAuthnName() string {
	return u.Username
}

func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	fmt.Printf("WebAuthnCredentials called, returning %d credentials", len(u.Credentials))
	return u.Credentials
}

func GetUser(username string) *User {
	return &User{
		Id:          []byte(username),
		Username:    username,
		DisplayName: username,
		Credentials: []webauthn.Credential{},
	}
}
