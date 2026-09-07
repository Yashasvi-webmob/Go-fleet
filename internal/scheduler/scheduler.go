package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Yashasvi-webmob/gofleet/internal/executor"
	"github.com/Yashasvi-webmob/gofleet/internal/job"
	"github.com/google/uuid"
)

type Scheduler struct {
	jobs chan uuid.UUID
	svc  *job.Service
	wg   sync.WaitGroup
}

func New(svc *job.Service, queueSize, workerCount int) *Scheduler {
	return &Scheduler{jobs: make(chan uuid.UUID, queueSize), svc: svc}
}

func (s *Scheduler) Enqueue(id uuid.UUID) {
	s.jobs <- id
}

func (s *Scheduler) Start(ctx context.Context, workerCount int) {

	for i := range workerCount {
		s.wg.Add(1)
		go func(workerID int) {
			defer s.wg.Done()
			for jobID := range s.jobs {
				j, err := s.svc.GetJob(ctx, jobID)
				if err != nil {
					continue
				}
				var currentStatus job.Status = job.StatusPending
				s.svc.Transition(ctx, jobID, job.StatusRunning)
				timeoutSecs := time.Duration(j.TimeoutSeconds) * time.Second
				res, execErr := executor.Run(ctx, j.Command, timeoutSecs)
				if execErr != nil {
					switch {
					case errors.Is(execErr, context.DeadlineExceeded):
						s.svc.Transition(ctx, jobID, job.StatusTimeout)
						currentStatus = job.StatusTimeout
					case errors.Is(execErr, context.Canceled):
						s.svc.Transition(ctx, jobID, job.StatusCancelled)
						currentStatus = job.StatusCancelled
					default:
						s.svc.Transition(ctx, jobID, job.StatusFailed)
						currentStatus = job.StatusFailed
					}
				} else {
					s.svc.Transition(ctx, jobID, job.StatusSuccess)
					currentStatus = job.StatusSuccess
				}

				fmt.Printf("Job with ID: %v set to Status: %v, Response from executor: %v \n", jobID, currentStatus, res)
			}
		}(i)
	}

}

func (s *Scheduler) Wait() {
	s.wg.Wait()
}

func (s *Scheduler) Close() {
	close(s.jobs)
}
