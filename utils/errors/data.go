package errors

import "strings"

const (
	ErrNotFound modelError = "accounts: resource not found"
	ErrIDInvalid modelError = "accounts: ID provided was invalid"
	ErrPasswordIncorrect modelError = "accounts: incorrect password provided"
	ErrPasswordRequired modelError = "accounts: password is required"
	ErrPasswordTooShort modelError = "accounts: password must be at least 8 characters"
	ErrEmailRequired modelError = "accounts: email address is required"
	ErrEmailInvalid modelError = "accounts: email address is invalid"
	ErrEmailTaken modelError = "accounts: email address is taken"
	ErrRememberTokenTooShort modelError = "accounts: remember token must be at least 32 bytes"
	ErrRememberTokenRequired modelError = "accounts: remember token is required"
)

type modelError string

func (err modelError) Error() string {
	return string(err)
}

func (err modelError) Public() string {
	s := strings.Replace(err.Error(), "accounts: ", "", 1)
	splits := strings.Split(s, " ")
	splits[0] = strings.Title(splits[0])
	return strings.Join(splits, " ")
}
