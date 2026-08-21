package job

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("job not found")

type InMemoryRepository struct {
	mu   sync.Mutex
	jobs map[uuid.UUID]*Job
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{jobs: make(map[uuid.UUID]*Job)}
}

func (r *InMemoryRepository) CreateJob(ctx context.Context, j *Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[j.ID] = j
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	j, ok := r.jobs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return j, nil
}

func (r *InMemoryRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	j, ok := r.jobs[id]
	if !ok {
		return ErrNotFound
	}
	j.Status = status
	return nil

}

func (r *InMemoryRepository) UpdateRetries(ctx context.Context, id uuid.UUID, retries int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	j, ok := r.jobs[id]
	if !ok {
		return ErrNotFound
	}
	j.Retries = retries
	return nil
}

var _ Repository = (*InMemoryRepository)(nil)
