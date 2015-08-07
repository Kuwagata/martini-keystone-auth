package auth

import "github.com/yuemin-li/Playground"

type Token string

type TokenValidator interface {
	ValidateToken(token string) bool
}

type IdentityValidator struct {
	auth_function func(token string) int
}

func (i IdentityValidator) ValidateToken(token string) bool {
	token_status := i.auth_function(token)
	if token_status == 200 {
		return true
	} else {
		return false
	}
}

func GetIdentityValidator(auth_url string) TokenValidator {
	return IdentityValidator{
		auth_function: func(token string) int {
			return identity.Auth(token, auth_url)
		}}
}
