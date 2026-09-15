package transport

import (
	"context"
	"errors"

	pb "github.com/Yashasvi-webmob/gofleet/gen"
	"github.com/Yashasvi-webmob/gofleet/internal/job"
	"github.com/Yashasvi-webmob/gofleet/internal/scheduler"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type JobServer struct {
	svc   *job.Service
	sched *scheduler.Scheduler
	pb.UnimplementedJobServiceServer
}

func NewJobServer(svc *job.Service, sched *scheduler.Scheduler) *JobServer {
	return &JobServer{svc: svc, sched: sched}
}

func toProtoJob(j *job.Job) *pb.Job {
	return &pb.Job{
		Id:        j.ID.String(),
		Status:    protoStatusMapper(j),
		Priority:  j.Priority,
		Command:   j.Command,
		CreatedAt: timestamppb.New(j.CreatedAt),
		UpdatedAt: timestamppb.New(j.UpdatedAt),
	}
}

func protoStatusMapper(j *job.Job) pb.JobStatus {
	switch j.Status {
	case job.StatusPending:
		return pb.JobStatus_JOB_STATUS_PENDING
	case job.StatusRunning:
		return pb.JobStatus_JOB_STATUS_RUNNING
	case job.StatusSuccess:
		return pb.JobStatus_JOB_STATUS_SUCCESS
	case job.StatusFailed:
		return pb.JobStatus_JOB_STATUS_FAILED
	case job.StatusTimeout:
		return pb.JobStatus_JOB_STATUS_TIMEOUT
	case job.StatusCancelled:
		return pb.JobStatus_JOB_STATUS_CANCELLED
	default:
		return pb.JobStatus_JOB_STATUS_UNSPECIFIED
	}
}

var _ pb.JobServiceServer = (*JobServer)(nil)

func (jobServer *JobServer) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.Job, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid jobID")
	}
	j, err := jobServer.svc.GetJob(ctx, id)
	if errors.Is(err, job.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "job not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProtoJob(j), nil
}

func (jobServer *JobServer) CreateJob(ctx context.Context, req *pb.SubmitJobRequest) (*pb.Job, error) {
	j, err := jobServer.svc.Submit(ctx, req.Command, int(req.TimeoutSeconds), req.Priority)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	newJob := toProtoJob(j)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	jobServer.sched.Enqueue(j.ID)
	return newJob, nil
}

func (jobServer *JobServer) CancelJob(ctx context.Context, req *pb.CancelJobRequest) (*empty.Empty, error) {
	jobID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid jobID")
	}
	transitionErr := jobServer.svc.Transition(ctx, jobID, job.StatusCancelled)
	if transitionErr != nil {
		switch {
		case errors.Is(transitionErr, job.ErrInvalidTransition):
			return nil, status.Error(codes.InvalidArgument, transitionErr.Error())
		case errors.Is(transitionErr, job.ErrRetriesExhausted):
			return nil, status.Error(codes.InvalidArgument, transitionErr.Error())
		default:
			return nil, status.Error(codes.Internal, transitionErr.Error())
		}

	}
	return nil, nil

}
