package apperrs

import (
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("email or password incorrect")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")

	ErrInternal     = errors.New("internal error")
	ErrInvalidToken = errors.New("invalid authorization token")
	ErrEmptyToken   = errors.New("auth token is empty, please login first")

	ErrAssetNotFound = errors.New("asset not found")
)
