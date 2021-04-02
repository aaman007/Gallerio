package models

import "strings"

const (
	ErrNotFound          modelError = "models: resource not found"
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	ErrPasswordRequired  modelError = "models: password is required"
	ErrPasswordTooShort  modelError = "models: password must be at least 8 characters"
	ErrEmailRequired     modelError = "models: email address is required"
	ErrEmailInvalid      modelError = "models: email address is invalid"
	ErrEmailTaken        modelError = "models: email address is taken"
	ErrTitleRequired     modelError = "models: title is required"
	ErrTokenInvalid      modelError = "models: token is invalid"
	ErrProviderRequired  modelError = "models: provider is required"
	
	ErrIDInvalid             privateError = "models: ID provided was invalid"
	ErrRememberTokenTooShort privateError = "models: remember token must be at least 32 bytes"
	ErrRememberTokenRequired privateError = "models: remember token is required"
	ErrUserIDRequired        privateError = "models: user ID was not provided"
)

type modelError string

func (err modelError) Error() string {
	return string(err)
}

func (err modelError) Public() string {
	s := strings.Replace(err.Error(), "models: ", "", 1)
	splits := strings.Split(s, " ")
	splits[0] = strings.Title(splits[0])
	return strings.Join(splits, " ")
}

type privateError string

func (err privateError) Error() string {
	return string(err)
}
