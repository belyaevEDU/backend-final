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
	Submit(userID string) (string, error)
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
