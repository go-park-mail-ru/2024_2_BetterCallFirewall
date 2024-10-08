package myErr

import (
	"errors"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrNoAuth               = errors.New("no auth")
	ErrWrongEmailOrPassword = errors.New("wrong email or password")
	ErrNonValidEmail        = errors.New("invalid email")
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionAlreadyExists = errors.New("session already exists")
)
