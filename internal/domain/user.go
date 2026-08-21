package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
)

type User struct {
	ID    string
	Login string
	// passwords are not stored in plaintext, repository will hold a bcrypt hash
	Password string
}
