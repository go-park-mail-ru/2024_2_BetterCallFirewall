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
	ErrToLargeFile          = errors.New("file too large")
	ErrNoMoreContent        = errors.New("no more content")
	ErrWrongFiletype        = errors.New("wrong type of file")
	ErrPostNotFound         = errors.New("post not found")
	ErrAccessDenied         = errors.New("access denied")
	ErrInternal             = errors.New("internal error")
	ErrWrongOwner           = errors.New("wrong owner")
	ErrSameUser             = errors.New("same user")
	ErrEmptyId              = errors.New("empty id")
	ErrBigId                = errors.New("id is too big")
	ErrProfileNotFound      = errors.New("profile not found")
	ErrAnotherService       = errors.New("another service")
	ErrInvalidQuery         = errors.New("invalid query parameter")
	ErrInvalidContext       = errors.New("invalid context parameter")
	ErrWrongDateFormat      = errors.New("wrong date format")
)
