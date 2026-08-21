package port

import (
	"github.com/belyaevedu/remote-code-service/internal/domain"
)

type TaskRepository interface {
	Save(task *domain.Task) error
	Get(id string) (*domain.Task, error)
	UpdateStatus(id string, status domain.TaskStatus) error
	SaveResult(id string, result *domain.Result) error
}

type TaskService interface {
	Submit() (string, error)
	Status(id string) (domain.TaskStatus, error)
	Result(id string) (*domain.Result, error)
}
