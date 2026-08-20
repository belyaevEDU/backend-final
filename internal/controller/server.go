package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Api struct {
	httpServer *http.Server
	shutdown   time.Duration
}

func NewApi(address string, router *chi.Mux, shutdownTimeout time.Duration) *Api {
	if shutdownTimeout <= 0 {
		shutdownTimeout = 10 * time.Second
	}
	return &Api{
		httpServer: &http.Server{
			Addr:              address,
			Handler:           router,
			ReadHeaderTimeout: 10 * time.Second, // gosec yelled at me the last time i didnt have this
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
		},
		shutdown: shutdownTimeout,
	}
}

func (a *Api) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	log.Printf("http server listening on %s", a.httpServer.Addr)

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Printf("shutdown signal received, draining connections (timeout %s)", a.shutdown)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.shutdown)
	defer cancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v\n", err)
		return err
	}
	return nil
}
