package handlers

import (
	"errors"
	"net/http"

	"github.com/belyaevedu/remote-code-service/internal/domain"
	"github.com/belyaevedu/remote-code-service/internal/port"

	"github.com/go-chi/chi/v5"
)

const (
	messageMissingAuthedUser = "missing authenticated user"
)

type TaskHandlers struct {
	taskSvc port.TaskService
}

func NewTaskHandlers(taskSvc port.TaskService) *TaskHandlers {
	return &TaskHandlers{taskSvc: taskSvc}
}

// POST /task
type taskCreateResponse struct {
	TaskID string `json:"task_id"`
}

// GET /status/{task_id}
type taskStatusResponse struct {
	Status string `json:"status"`
}

// GET /result/{task_id}
type taskResultResponse struct {
	Result domain.Result `json:"result"`
}

// @Summary Create a task
// @Description Creating a task owned by an authenticated user
// @Tags task
// @Accept json
// @Produce json
// @Success 201 {object} taskCreateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /task [post]
func (h *TaskHandlers) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		unauthorizedResponseHelper(w, messageMissingAuthedUser)
		return
	}

	id, err := h.taskSvc.Submit(userID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	WriteJSON(w, http.StatusCreated, taskCreateResponse{TaskID: id})
}

// @Summary Get the status of a certain task
// @Description Getting the status of a task by id. Only the owner of the task can access it
// @Tags task
// @Accept json
// @Produce json
// @Param task_id path string true "task id"
// @Success 200 {object} taskStatusResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /status/{task_id} [get]
func (h *TaskHandlers) Status(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		unauthorizedResponseHelper(w, messageMissingAuthedUser)
		return
	}

	id := chi.URLParam(r, "task_id")

	status, err := h.taskSvc.Status(userID, id)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, taskStatusResponse{Status: string(status)})
}

// @Summary Get the result of a certain task
// @Description Getting the result of a task by id. Only the owner of the task can access it
// @Tags task
// @Accept json
// @Produce json
// @Param task_id path string true "task id"
// @Success 200 {object} taskResultResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /result/{task_id} [get]
func (h *TaskHandlers) Result(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		unauthorizedResponseHelper(w, messageMissingAuthedUser)
		return
	}

	id := chi.URLParam(r, "task_id")

	result, err := h.taskSvc.Result(userID, id)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, taskResultResponse{Result: *result})
}

func userIDFromContext(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(UserIDKey).(string)
	return id, ok && id != ""
}

// both unknown tasks and tasks owned by another user are reported as 404 with the same body
func writeTaskError(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrTaskNotFound) || errors.Is(err, domain.ErrAccessDenied) {
		WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: domain.ErrTaskNotFound.Error()})
		return
	}
	WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
}
