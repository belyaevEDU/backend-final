package domain

import "errors"

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	UserID    string
	SessionID string
}
