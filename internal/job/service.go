package job

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

const MaxRetries = 3

var allowedTransition = map[Status][]Status{
	StatusPending: {StatusRunning, StatusCancelled},
	StatusRunning: {StatusSuccess, StatusTimeout, StatusFailed, StatusCancelled},
	StatusFailed:  {StatusPending},
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ? Using the slice in built function
/*
func canTransition(from Status, to Status) bool {
	for _, allowed := range allowedTransition[from] {
		if allowed == to {
			return true
		}
	}
	return false
}*/

func (s *Service) Submit(ctx context.Context, command string, timeout int, priority string) (*Job, error) {

	var newJob = &Job{
		ID:             uuid.New(),
		Status:         StatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Command:        command,
		TimeoutSeconds: timeout,
		Priority:       priority,
		Retries:        0,
	}
	err := s.repo.CreateJob(ctx, newJob)
	if err != nil {
		return nil, fmt.Errorf("submit job: %w", err)
	}
	return newJob, nil

}

func (s *Service) GetJob(ctx context.Context, id uuid.UUID) (*Job, error) {
	j, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job status: %w", err)
	}
	return j, nil
}

func (s *Service) Transition(ctx context.Context, id uuid.UUID, status Status) error {
	currentJob, statusErr := s.GetJob(ctx, id)
	if statusErr != nil {
		return fmt.Errorf("transition get job status: %w", statusErr)
	}

	currentJobStatus := currentJob.Status
	var transitionVal = allowedTransition[currentJobStatus]

	if ok := slices.Contains(transitionVal, status); ok == false {
		return fmt.Errorf("status not in allowed transitions; currentJobStatus: %s , status given: %s", currentJobStatus, status)
	}
	// if acceptable := canTransition(currentJobStatus, status); acceptable == false {
	// 	return fmt.Errorf("status not in allowed transitions: %v", acceptable)
	// }

	// * Retries logic
	if currentJobStatus == StatusFailed && status == StatusPending {

		retriesDone := currentJob.Retries
		if retriesDone >= MaxRetries {
			return fmt.Errorf("no more retries left for this job: %v", currentJob.ID)
		} else {
			if err := s.repo.UpdateRetries(ctx, id, retriesDone+1); err != nil {
				return fmt.Errorf("updating retires: %w", err)
			}
		}
	}

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("job status transition: %w", err)
	}
	return nil
}
