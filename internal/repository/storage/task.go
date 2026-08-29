package storage

import (
	"github.com/belyaevedu/remote-code-service/internal/domain"
)

func (r *Repository) Save(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task

	return nil
}

func (r *Repository) Get(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tasks[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}

	// returning a clone so the caller sees an unchanging snapshot
	clone := *t
	return &clone, nil
}

func (r *Repository) UpdateStatus(id string, status domain.TaskStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tasks[id]
	if !ok {
		return domain.ErrTaskNotFound
	}

	t.Status = status

	return nil
}

func (r *Repository) SaveResult(id string, result *domain.Result) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tasks[id]
	if !ok {
		return domain.ErrTaskNotFound
	}

	t.Status = domain.StatusReady
	t.Result = result

	return nil
}
