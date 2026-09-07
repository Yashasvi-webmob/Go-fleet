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

## V2 — not started
- protobuf, gRPC `JobService`/`WorkerService`, generated code, CLI-over-gRPC

## V3 — not started
- Bidirectional streaming, live logs, heartbeats

## V4 — not started
- Postgres, Docker Compose, CI/CD, metrics, integration tests

## Notes
- Git history is currently thin relative to the actual work done (one commit so far) — worth committing per-milestone going forward so the repo's history reflects the real build process.
