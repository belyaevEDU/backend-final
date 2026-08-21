package repository

import (
	"github.com/belyaevedu/remote-code-service/internal/domain"
)

func (r *Repository) SaveUser(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clone := *user
	r.users[clone.Login] = &clone

	return nil
}

func (r *Repository) GetUserByLogin(login string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[login]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	clone := *u
	return &clone, nil
}

func (r *Repository) GetUserByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.ID == id {
			clone := *u
			return &clone, nil
		}
	}

	return nil, domain.ErrUserNotFound
}
