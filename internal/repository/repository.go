package repository

import (
	"sync"

	"github.com/belyaevedu/remote-code-service/internal/domain"
	"github.com/belyaevedu/remote-code-service/internal/port"
)

type Repository struct {
	mu sync.RWMutex

	tasks    map[string]*domain.Task
	users    map[string]*domain.User    // key - login
	sessions map[string]*domain.Session // key - session id
}

// compile-time asserts
var (
	_ port.TaskRepository    = (*Repository)(nil)
	_ port.UserRepository    = (*Repository)(nil)
	_ port.SessionRepository = (*Repository)(nil)
)

func New() *Repository {
	return &Repository{
		tasks:    make(map[string]*domain.Task),
		users:    make(map[string]*domain.User),
		sessions: make(map[string]*domain.Session),
	}
}
