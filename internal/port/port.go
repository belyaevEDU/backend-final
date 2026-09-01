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
	Submit(ctx context.Context, userID string, sub domain.Submission) (string, error)
	Status(userID, id string) (domain.TaskStatus, error)
	Result(userID, id string) (*domain.Result, error)
}

type UserRepository interface {
	SaveUser(user *domain.User) error
	GetUserByID(id string) (*domain.User, error)
	GetUserByLogin(login string) (*domain.User, error)
}

type UserService interface {
	Register(login, password string) error
	Login(login, password string) (string, error)
}

type SessionRepository interface {
	CreateSession(session *domain.Session) error
	GetSession(sessionID string) (*domain.Session, error)
	DeleteSession(sessionID string) error
}

type AuthService interface {
	Authenticate(token string) (string, error)
}

type CodeExecutor interface {
	Execute(ctx context.Context, msg domain.TaskMessage) (domain.ExecutionResult, error)
}

type TaskPublisher interface {
	Publish(ctx context.Context, msg domain.TaskMessage) error
}

// processes a single task message consumed from the queue
// returning an error requeues the message once. a message that already failed once is dropped
type TaskHandler func(ctx context.Context, msg domain.TaskMessage) error

type TaskConsumer interface {
	Consume(ctx context.Context, handler TaskHandler) error
}
