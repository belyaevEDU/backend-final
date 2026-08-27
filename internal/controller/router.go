package controller

import (
	"github.com/go-chi/chi/v5"

	"github.com/belyaevedu/remote-code-service/internal/controller/handlers"
	"github.com/belyaevedu/remote-code-service/internal/port"
)

func NewRouter(t *handlers.TaskHandlers, u *handlers.UserHandlers, auth port.AuthService) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", u.Register)
	r.Post("/login", u.Login)

	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware(auth))
		r.Post("/task", t.Create)
		r.Get("/status/{task_id}", t.Status)
		r.Get("/result/{task_id}", t.Result)
	})

	return r
}
