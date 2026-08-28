package domain

import (
	"errors"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrAccessDenied = errors.New("access denied")
)

type TaskStatus string

const (
	StatusInProgress TaskStatus = "in_progress"
	StatusReady      TaskStatus = "ready"
)

type Result struct {
	Output string `json:"output,omitempty"`
}

type Task struct {
	ID     string
	UserID string
	Status TaskStatus
	Result *Result
}
