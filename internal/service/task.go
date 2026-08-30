package service

import (
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

func NewTaskService(repo port.TaskRepository, processingTime time.Duration) *TaskService {
	if processingTime <= 0 {
		processingTime = 2 * time.Second
	}
	return &TaskService{
		repo:           repo,
		processingTime: processingTime,
	}
}

func (s *TaskService) Submit(userID string) (string, error) {
	if userID == "" {
		return "", domain.ErrAccessDenied
	}

	id := uuid.NewString()

	task := &domain.Task{
		ID:     id,
		UserID: userID,
		Status: domain.StatusInProgress,
	}

	if err := s.repo.Save(task); err != nil {
		return "", err
	}

	go s.process(id)

	return id, nil
}

func (s *TaskService) Status(userID, id string) (domain.TaskStatus, error) {
	t, err := s.getOwnedTask(userID, id)
	if err != nil {
		return "", err
	}
	return t.Status, nil
}

func (s *TaskService) Result(userID, id string) (*domain.Result, error) {
	t, err := s.getOwnedTask(userID, id)
	if err != nil {
		return nil, err
	}
	if t.Result == nil {
		return nil, domain.ErrTaskNotFound
	}
	return t.Result, nil
}

// fetches the task by id and verifies it belongs to a set user
// task owned by someone else results in ErrAccessDenied
func (s *TaskService) getOwnedTask(userID, id string) (*domain.Task, error) {
	if userID == "" {
		return nil, domain.ErrAccessDenied
	}

	t, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	if t.UserID != userID {
		return nil, domain.ErrAccessDenied
	}

	return t, nil
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
