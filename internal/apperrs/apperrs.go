package apperrs

import (
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("email or password incorrect")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInternal           = errors.New("internal error")
	ErrInvalidToken       = errors.New("invalid authorization token")
)
