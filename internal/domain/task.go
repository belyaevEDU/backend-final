package domain

import (
	"errors"
)

var (
	ErrTaskNotFound          = errors.New("task not found")
	ErrAccessDenied          = errors.New("access denied")
	ErrUnsupportedTranslator = errors.New("unsupported translator")
	ErrExecutionTimeout      = errors.New("execution timed out")
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

type TaskMessage struct {
	TaskID     string `json:"task_id"`
	Translator string `json:"translator"`
	Code       string `json:"code"`
}

type ExecutionRequest struct {
	Name       string // the executor-side task name, expected to be unique
	Translator string
	Code       string
}

type ExecutionResult struct {
	Output   string
	ExitCode int // -1 when unavailable
	Failed   bool
}
