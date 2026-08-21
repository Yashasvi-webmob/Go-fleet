package job_test

import (
	"context"
	"testing"
	"time"

	"github.com/Yashasvi-webmob/gofleet/internal/job"
	"github.com/google/uuid"
)

func newTestService() (*job.Service, job.Repository) {
	repo := job.NewInMemoryRepository()
	return job.NewService(repo), repo
}

func TestServiceSubmit(t *testing.T) {
	svc, _ := newTestService()
	ctx := context.Background()
	got, err := svc.Submit(ctx, "working check", 20, "high")
	if err != nil {
		t.Fatalf("Submit() err: %v", err)
	}
	if got.Status != job.StatusPending {
		t.Fatalf("Status initial not pending, current:%v", got.Status)
	}
	if got.Retries != 0 {
		t.Fatalf("Retries initial not 0, current:%d", got.Retries)
	}
	if got.ID == uuid.Nil {
		t.Error("ID not generated")
	}

}

func TestService_Transition(t *testing.T) { // ? explain working in detail
	cases := []struct {
		name    string
		from    job.Status
		to      job.Status
		wantErr bool
	}{
		{"pending to running", job.StatusPending, job.StatusRunning, false},
		{"pending to failed", job.StatusPending, job.StatusFailed, true},
		{"running to cancelled", job.StatusRunning, job.StatusCancelled, false},
		{"pending to cancelled", job.StatusPending, job.StatusCancelled, false},
		{"cancelled to pending", job.StatusCancelled, job.StatusPending, true},
		{"pending to timeout", job.StatusPending, job.StatusTimeout, true},
		{"timeout to pending", job.StatusTimeout, job.StatusPending, true},
		{"pending to success", job.StatusPending, job.StatusSuccess, true},
		{"success to timeout", job.StatusSuccess, job.StatusTimeout, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, repo := newTestService()
			ctx := context.Background()
			seed := job.Job{ID: uuid.New(), Status: tc.from, Priority: "high", Retries: 0, Command: "sample", TimeoutSeconds: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			if err := repo.CreateJob(ctx, &seed); err != nil {
				t.Fatalf("seed CreateJob() error: %v", err)
			}
			err := svc.Transition(ctx, seed.ID, tc.to)
			if (err != nil) != tc.wantErr {
				t.Errorf("Transition(%v -> %v) error :%v wantErr: %v", tc.from, tc.to, err, tc.wantErr)
			}
		})
	}
}

func TestService_Retries(t *testing.T) { // ? explain working in detail
	cases := []struct {
		name          string
		retriesUpdate int
		wantErr       bool
	}{
		{"Retries: 1", 1, false},
		{"Retries: 0", 0, false},
		{"Retries: 2", 2, false},
		{"Retries: 3", 3, true},
		{"Retries: 5", 5, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, repo := newTestService()
			ctx := context.Background()
			seed := job.Job{ID: uuid.New(), Status: job.StatusFailed, Priority: "high", Retries: tc.retriesUpdate, Command: "sample", TimeoutSeconds: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			if err := repo.CreateJob(ctx, &seed); err != nil {
				t.Fatalf("CreatedJob() error %v", err)
			}

			err := svc.Transition(ctx, seed.ID, job.StatusPending)
			if (err != nil) != tc.wantErr {
				t.Fatalf(" Transition() err, wantErr (retries exhausted), wantErr: %v", tc.wantErr)
			}
		})
	}
}
