package handlers

import (
	"errors"
	"net/http"

	"github.com/belyaevedu/remote-code-service/internal/domain"
	"github.com/belyaevedu/remote-code-service/internal/port"

	"github.com/go-chi/chi/v5"
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
// @Description Creating a task
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
	id, err := h.taskSvc.Submit()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	WriteJSON(w, http.StatusCreated, taskCreateResponse{TaskID: id})
}

// @Summary Get the status of a certain task
// @Description Getting the status of a task by id
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
	id := chi.URLParam(r, "task_id")

	status, err := h.taskSvc.Status(id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: domain.ErrTaskNotFound.Error()})
			return
		}
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	WriteJSON(w, http.StatusOK, taskStatusResponse{Status: string(status)})
}

// @Summary Get the result of a certain task
// @Description Getting the result of a task by id
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
	id := chi.URLParam(r, "task_id")

	result, err := h.taskSvc.Result(id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: domain.ErrTaskNotFound.Error()})
			return
		}
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	WriteJSON(w, http.StatusOK, taskResultResponse{Result: *result})
}
