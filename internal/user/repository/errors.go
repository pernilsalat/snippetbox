package repository

import "errors"

var (
	ErrInvalidCredentials = errors.New("models: invalid credentials provided")
	ErrDuplicateEmail     = errors.New("models: duplicate email address provided")
)
