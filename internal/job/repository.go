package job

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateJob(ctx context.Context, j *Job) error
	GetByID(ctx context.Context, id uuid.UUID) (*Job, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error
	UpdateRetries(ctx context.Context, id uuid.UUID, retries int) error
}
