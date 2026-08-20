package port

import (
	"context"

	"github.com/belyaevedu/remote-code-service/internal/domain"
)

type TaskRepository interface {
	Save(task *domain.Task) error
	Get(id string) (*domain.Task, error)
	UpdateStatus(id string, status domain.TaskStatus) error
	SaveResult(id string, result *domain.Result) error
}

type TaskService interface {
	Submit(ctx context.Context) (string, error)
	Status(ctx context.Context, id string) (domain.TaskStatus, error)
	Result(ctx context.Context, id string) (*domain.Result, error)
}
