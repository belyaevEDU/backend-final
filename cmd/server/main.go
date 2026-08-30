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

// @title HTTP REST API for a safe remote code execution service
// @version 1.0
// @description The foundation for a safe remote code execution service
//
// @contact.name Veniamin Belyaev
// @contact.email veniamin@belyaev.work
// @license.name MIT
//
// @BasePath /
//
// @accept json
// @produce json
// @schemes http
//
// @securitydefinitions.bearerauth BearerAuth
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	repo := repository.New()

	taskService := service.NewTaskService(repo, cfg.ProcessingTime)
	userService := service.NewUserService(repo, repo)

	taskHandler := handlers.NewTaskHandlers(taskService)
	userHandler := handlers.NewUserHandlers(userService)

	router := controller.NewRouter(taskHandler, userHandler, userService)
	server := controller.NewApi(cfg.HTTPAddr, router, cfg.ShutdownTimeout)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := server.Start(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
