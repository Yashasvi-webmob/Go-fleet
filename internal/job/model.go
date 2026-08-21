package job

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	StatusPending Status = iota
	StatusRunning
	StatusSuccess
	StatusFailed
	StatusTimeout
	StatusCancelled
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "Status Pending"
	case StatusRunning:
		return "Status Running"
	case StatusSuccess:
		return "Status Success"
	case StatusCancelled:
		return "Status Cancelled"
	case StatusFailed:
		return "Status Failed"
	case StatusTimeout:
		return "Status Timed out"
	default:
		return "Invalid Status"
	}
}

type Job struct {
	ID             uuid.UUID `json:"id"`
	Status         Status    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Command        string    `json:"command"`
	TimeoutSeconds int       `json:"timeout_seconds"`
	Priority       string    `json:"priority"`
	Retries        int       `json:"retires"`
}
