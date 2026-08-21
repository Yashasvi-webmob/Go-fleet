package job_test

import (
	"context"
	"errors"
	"github.com/Yashasvi-webmob/gofleet/internal/job"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestRepo_CreateAndGet(t *testing.T) {
	InMemoryRepository := job.NewInMemoryRepository()
	ctx := context.Background() // ? *Claude* Explain this, this was not fully understood
	jobBody := &job.Job{ID: uuid.New(), Status: job.StatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now(), Command: "Sample Command", TimeoutSeconds: 10, Priority: "High", Retries: 0}
	if err := InMemoryRepository.CreateJob(ctx, jobBody); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := InMemoryRepository.GetByID(ctx, jobBody.ID)
	if err != nil {
		t.Fatalf("GetById() error %v", err)
	}

	if got.Status != jobBody.Status {
		t.Errorf("Status mismatched given= %q,got = %q", jobBody.Status, got.Status)
	}
}

func TestRepo_GetByInvalidId(t *testing.T) {
	InMemoryRepository := job.NewInMemoryRepository()
	var invalidId = uuid.New()
	ctx := context.Background()
	_, err := InMemoryRepository.GetByID(ctx, invalidId)
	if !errors.Is(err, job.ErrNotFound) {
		t.Fatalf("GetById() error = %v, want ErrNotFound", err)
	}
}

func TestRepo_UpdateStatus(t *testing.T) {
	InMemoryRepository := job.NewInMemoryRepository()
	ctx := context.Background()
	jobBody := &job.Job{ID: uuid.New(), Status: job.StatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now(), Command: "Sample Command", TimeoutSeconds: 10, Priority: "High", Retries: 0}
	if err := InMemoryRepository.CreateJob(ctx, jobBody); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	var updatedStatus = job.StatusFailed
	if updateErr := InMemoryRepository.UpdateStatus(ctx, jobBody.ID, updatedStatus); updateErr != nil {
		t.Fatalf("UpdateStatus() error %v", updateErr)
	}
	j, err := InMemoryRepository.GetByID(ctx, jobBody.ID)
	if err != nil {
		t.Fatalf("GetById() error %v", err)
	}
	if j.Status != updatedStatus {
		t.Errorf("Test UpdateStatus() status different: Set status:%v Incoming Status:%v", updatedStatus, j.Status)
	}
}

func TestRepo_UpdateRetries(t *testing.T) {
	InMemoryRepository := job.NewInMemoryRepository()
	ctx := context.Background()
	jobBody := &job.Job{ID: uuid.New(), Status: job.StatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now(), Command: "Sample Command", TimeoutSeconds: 10, Priority: "High", Retries: 0}
	if err := InMemoryRepository.CreateJob(ctx, jobBody); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	var updateRetries = 2
	if updateErr := InMemoryRepository.UpdateRetries(ctx, jobBody.ID, updateRetries); updateErr != nil {
		t.Fatalf("UpdateStatus() error %v", updateErr)
	}

	j, err := InMemoryRepository.GetByID(ctx, jobBody.ID)
	if err != nil {
		t.Fatalf("GetById() error %v", err)
	}
	if j.Retries != updateRetries {
		t.Errorf("Test UpdateRetries() retries different: Set retries:%v Incoming retries:%v", updateRetries, j.Retries)
	}
}

func TestRepo_NotFoundId(t *testing.T) {
	repo := job.NewInMemoryRepository()
	ctx := context.Background()
	randomID := uuid.New()
	cases := []struct {
		name string
		run  func() error
	}{
		{"GetById", func() error {
			_, err := repo.GetByID(ctx, randomID)
			return err
		}},
		{"UpdateStatus", func() error {
			err := repo.UpdateStatus(ctx, randomID, job.StatusTimeout)
			return err
		}},
		{"UpdateRetries", func() error {
			err := repo.UpdateRetries(ctx, randomID, 5)
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.run(); !errors.Is(err, job.ErrNotFound) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}
