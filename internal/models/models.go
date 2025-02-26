package models

import "github.com/go-webauthn/webauthn/webauthn"

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
