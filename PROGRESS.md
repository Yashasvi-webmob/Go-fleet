# GoFleet — Progress Log

Minimal phase tracker. Updated as milestones complete.

## Phase 0 — Scaffolding
- [x] `go.mod` + directory skeleton

## V1 — in-process, no gRPC yet
- [x] M1 — Job domain model (`internal/job/model.go`): `Job`, `Status` enum, `Stringer`
- [x] M2 — `Repository` interface + `InMemoryRepository` (mutex-guarded map)
- [x] M3 — `Service`: FSM transitions, retries, error wrapping
- [x] M4 — Table-driven unit tests (repository + service)
- [x] M5 — `executor.Run`: `os/exec` + context timeout, `context.DeadlineExceeded` detection
- [x] M6 — `cmd/cli`: `flag`-based subcommands (`submit`, `status`, `cancel`), manual DI
- [x] M7 — `internal/scheduler`: goroutine worker pool, channel-based job queue, full Submit → Run → Transition pipeline (verified end-to-end: success/fail/timeout all land correctly)
- [x] M8 — Graceful shutdown (`cmd/server`: `signal.NotifyContext`, block on `ctx.Done()`, `Scheduler.Close()` + `Wait()`) — verified live: SIGINT drains workers and exits cleanly
- [x] M9 — Review pass: full build/vet/fmt/test sweep clean; dead code + stray dev comments removed

**V1 complete.** Known, accepted limitations carried into V2+ rather than fixed now:
- `internal/scheduler` has no unit tests yet — `executor.Run` is called directly (not via an injectable interface), so testing it deterministically needs that refactor first. Revisit alongside V2.
- "Worker registration" (plan.md's V1 criteria) is satisfied loosely — a fixed in-process worker pool, not a real registration protocol. Becomes literal with `WorkerService.Connect` in V2/V3.

## V2 — gRPC (JobService only)
- [x] M10 — Toolchain: `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc` installed
- [x] M11 — `api/proto/gofleet.proto`: `Job`/`JobStatus`, `SubmitJobRequest`/`GetJobRequest`/`CancelJobRequest`, `JobService` (SubmitJob→`CreateJob` rpc, GetJob, CancelJob)
- [x] M12 — Generated code (`gen/gofleet.pb.go`, `gen/gofleet_grpc.pb.go`) reviewed: messages, nil-safe getters, `JobServiceServer`/`JobServiceClient` interfaces
- [x] M13 — `internal/transport/grpc/job_server.go`: `JobServer` implements `JobServiceServer`, wraps `job.Service` + `scheduler.Scheduler`; explicit status enum mapping (not a numeric cast); `job.ErrInvalidTransition`/`ErrRetriesExhausted` sentinels mapped to `codes.FailedPrecondition`-style errors
- [x] M14 — `cmd/server`: real `grpc.NewServer()` + `net.Listen` + `GracefulStop()` alongside scheduler shutdown; `cmd/cli`: full `JobServiceClient` conversion (no more in-process `job.Service`) — verified live end-to-end (submit → execute → status → correctly-rejected cancel → graceful shutdown)

**V2 (JobService) complete.** Deliberately deferred, per plan.md's own V2/V3 split:
- `WorkerService` (`Connect`, bidirectional streaming) — real content is V3's (worker registration, streaming job assignment), so it wasn't stubbed out early.
- `WatchLogs` — same reasoning, needs V3's streaming groundwork.
- Minor cleanup pending: a few `gofmt -w .`-fixable whitespace nits.

## V3 — not started
- Bidirectional streaming, live logs, heartbeats, `WorkerService`

## V4 — not started
- Postgres, Docker Compose, CI/CD, metrics, integration tests

## Notes
- Git history is currently thin relative to the actual work done — commit per-milestone going forward so the repo's history reflects the real build process.
