package service

import (
	"context"
	"log"
	"time"

	"github.com/belyaevedu/remote-code-service/internal/domain"
	"github.com/belyaevedu/remote-code-service/internal/port"

	"github.com/google/uuid"
)

const (
	outputMessage = "puk"
)

type TaskService struct {
	repo           port.TaskRepository
	processingTime time.Duration
}

// compile-time assert that task's TaskService struct
// implements port's TaskService interface
var _ port.TaskService = (*TaskService)(nil)

func New(repo port.TaskRepository, processingTime time.Duration) *TaskService {
	if processingTime <= 0 {
		processingTime = 2 * time.Second
	}
	return &TaskService{
		repo:           repo,
		processingTime: processingTime,
	}
}

func (s *TaskService) Submit(ctx context.Context) (string, error) {
	id := uuid.NewString()

	task := &domain.Task{
		ID:     id,
		Status: domain.StatusInProgress,
	}

	if err := s.repo.Save(task); err != nil {
		return "", err
	}

	go s.process(id)

	return id, nil
}

func (s *TaskService) Status(ctx context.Context, id string) (domain.TaskStatus, error) {
	t, err := s.repo.Get(id)
	if err != nil {
		return "", err
	}
	return t.Status, nil
}

func (s *TaskService) Result(ctx context.Context, id string) (*domain.Result, error) {
	t, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	if t.Result == nil {
		return nil, domain.ErrTaskNotFound
	}
	return t.Result, nil
}

func (s *TaskService) process(id string) {
	time.Sleep(s.processingTime)

	result := &domain.Result{
		Output: outputMessage,
	}

	if err := s.repo.SaveResult(id, result); err != nil {
		log.Printf("Failed to save task result %s: %v\n", id, err)
		return
	}
	log.Printf("Task finished: %s\n", id)
}
