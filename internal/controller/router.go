package controller

import (
	"github.com/belyaevedu/remote-code-service/internal/controller/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.TaskHandlers) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/task", h.Create)
	r.Get("/status/{task_id}", h.Status)
	r.Get("/result/{task_id}", h.Result)

	return r
}
