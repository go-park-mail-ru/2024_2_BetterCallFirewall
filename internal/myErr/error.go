package myErr

import (
	"errors"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrNoAuth               = errors.New("no session found")
	ErrWrongEmailOrPassword = errors.New("wrong email or password")
	ErrNonValidEmail        = errors.New("invalid email")
)
