package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/belyaevedu/remote-code-service/internal/config"
	"github.com/belyaevedu/remote-code-service/internal/controller"
	"github.com/belyaevedu/remote-code-service/internal/controller/handlers"
	"github.com/belyaevedu/remote-code-service/internal/repository"
	"github.com/belyaevedu/remote-code-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	repo := repository.New()

	taskService := service.New(repo, cfg.ProcessingTime)
	userService := service.NewUserService(repo, repo)

	taskHandler := handlers.New(taskService)
	userHandler := handlers.NewUserHandlers(userService)

	router := controller.NewRouter(taskHandler, userHandler, userService)
	server := controller.NewApi(cfg.HTTPAddr, router, cfg.ShutdownTimeout)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	if err := server.Start(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
