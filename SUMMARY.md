# GoFleet — Project Summary & Self-Assessment

## What this repo is

A distributed job execution platform in Go — users submit shell-command jobs, a scheduler assigns them to workers, workers execute and report status. Modeled loosely on GitHub Actions runners / Celery / Kubernetes Jobs, built specifically to learn production Go patterns (not as a CRUD-app tutorial). Full spec in `plan.md`.

The build follows a strict rule: every line of implementation code was written by the repo owner, not generated. Guidance came as explained concepts, syntax snippets, and code review — never a finished solution. That's the whole point: fluency, not a finished artifact.

## Where it stands (see `PROGRESS.md` for the live log)

**V1 is complete**, and **V2's `JobService` slice is complete.** V1: domain model, repository, service layer with a hand-validated FSM, unit tests, a real `os/exec`-based executor with proper timeout/cancellation handling, a goroutine/channel worker pool, and a `cmd/server` binary with real graceful shutdown. V2: a hand-written `.proto` schema, generated protobuf/gRPC code, a `JobServer` implementation wrapping the existing `job.Service` + `Scheduler`, and both `cmd/server` (real `grpc.NewServer()`) and `cmd/cli` (now a genuine network client, no in-process business logic left) converted — verified live end-to-end: submit a job over gRPC, watch it actually execute and transition status, and watch an invalid cancel get correctly rejected with a semantically-appropriate gRPC error code. Known gaps: `internal/scheduler` still has no unit tests; `WorkerService`/streaming/`WatchLogs` are intentionally deferred to V3, since real worker registration only makes sense once streaming exists.

## What you should walk away knowing

- **Interfaces & structural typing** — small interfaces, satisfied implicitly, no `implements` keyword; why the `Repository` interface stays swappable for Postgres later.
- **Error handling idioms** — `%w` wrapping vs `%v`, `errors.Is`/sentinel errors, why the *thing you wrap* has to be the thing a caller will later check for.
- **Concurrency primitives** — goroutines, channels as a queue, the worker-pool pattern (shared channel, N consumers, no manual assignment logic needed), `sync.Mutex`/`RWMutex`, `sync.WaitGroup`, and the classic footguns: copying a struct that holds a mutex, closures capturing loop variables, `:=` shadowing inside nested blocks instead of reusing an outer variable.
- **`context.Context`** — as the root-context idiom, `WithTimeout`/cancellation propagation, and the non-obvious part: a killed subprocess's own error doesn't itself carry `context.DeadlineExceeded` — you have to check the context directly and deliberately wrap it if a caller needs to detect it later.
- **Manual dependency injection** — no framework; wiring happens by hand in `main()`, and every layer only depends on the interface it needs, not concrete types.
- **Table-driven testing** — plus the discipline of deriving expected values from the *spec*, not from the implementation being tested (a mistake that surfaced and got caught more than once here).
- **Package design** — organizing by domain (`job`, `scheduler`, `executor`, `transport/grpc`) instead of generic `utils`/`helpers`.
- **Protobuf/gRPC** — schema-first IDL, wire tags vs. field names, why enum zero-values matter, generated code as a boundary you call through but never hand-edit, the `UnimplementedXServer`-by-value forward-compatibility pattern, structured gRPC error codes (`codes.NotFound`, `codes.FailedPrecondition`) vs. plain Go errors.

Once V3-V4 land, add: bidirectional streaming, worker heartbeats/registration, Postgres persistence, Docker Compose, CI/CD.

## Is this showcase-worthy (Twitter/LinkedIn/GitHub) right now?

**Closer, but still not quite yet.** V1 + V2's `JobService` slice demonstrates solid Go fundamentals *and* a real typed service boundary over the network — a meaningfully bigger step than a lot of portfolio repos. It's not there yet mainly because there's no actual distributed worker (jobs still execute in the same process as the scheduler) and no persistence beyond memory.

What would actually make it worth a public post:
- **V3 (streaming + real `WorkerService`)** is the big remaining piece — this is what turns "a job queue with a gRPC API" into an actual distributed system with a separate worker process, live log streaming, and heartbeats.
- **V4 (Docker Compose + CI/CD + Postgres + integration tests)** is what makes it look production-grade rather than a weekend exercise.
- A proper root `README.md` with an architecture diagram, setup instructions, and example CLI usage — first impression matters as much as the code.
- More granular git history — commit per-milestone from here on so a viewer can see the actual iterative build process, one of the strongest signals the work is genuinely understood rather than AI-authored wholesale.

**The genuine differentiator you'll have, once you do post it:** you can actually explain every design decision in an interview — why the repository is an interface, why the FSM is a map instead of scattered `if`s, why `%w` matters, why the status enum needs an explicit mapping instead of a numeric cast, why `UnimplementedJobServiceServer` must be embedded by value. That's real and rare; most portfolio projects (AI-assisted or not) can't survive that kind of follow-up questioning. Worth saying explicitly in a post/README, not just implying it.

**Recommendation:** hold off on publicizing until at least V3 is done (ideally V4), keep committing incrementally from here on, and write the README last, once there's a finished story to tell.
