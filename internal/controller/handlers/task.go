package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/belyaevedu/remote-code-service/internal/domain"
	"github.com/belyaevedu/remote-code-service/internal/port"

	"github.com/go-chi/chi/v5"
)

type TaskHandlers struct {
	taskSvc port.TaskService
}

func New(taskSvc port.TaskService) *TaskHandlers {
	return &TaskHandlers{taskSvc: taskSvc}
}

// body returned by POST /task
type taskCreateResponse struct {
	TaskID string `json:"task_id"`
}

// body returned by GET /status/{task_id}
type taskStatusResponse struct {
	Status string `json:"status"`
}

// body returned by GET /result/{task_id}
type taskResultResponse struct {
	Result domain.Result `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// POST /task
func (h *TaskHandlers) Create(w http.ResponseWriter, r *http.Request) {
	id, err := h.taskSvc.Submit(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, taskCreateResponse{TaskID: id})
}

// GET /status/{task_id}
func (h *TaskHandlers) Status(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "task_id")

	status, err := h.taskSvc.Status(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: domain.ErrTaskNotFound.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, taskStatusResponse{Status: string(status)})
}

// GET /result/{task_id}
func (h *TaskHandlers) Result(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "task_id")

	result, err := h.taskSvc.Result(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: domain.ErrTaskNotFound.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, taskResultResponse{Result: *result})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("Error raised encoding a response into JSON: %v\n", err)
	}
}
