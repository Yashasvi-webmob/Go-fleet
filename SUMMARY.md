# GoFleet — Project Summary & Self-Assessment

## What this repo is

A distributed job execution platform in Go — users submit shell-command jobs, a scheduler assigns them to workers, workers execute and report status. Modeled loosely on GitHub Actions runners / Celery / Kubernetes Jobs, built specifically to learn production Go patterns (not as a CRUD-app tutorial). Full spec in `plan.md`.

The build follows a strict rule: every line of implementation code was written by the repo owner, not generated. Guidance came as explained concepts, syntax snippets, and code review — never a finished solution. That's the whole point: fluency, not a finished artifact.

## Where it stands (see `PROGRESS.md` for the live log)

**V1 is complete.** Domain model, repository, service layer with a hand-validated FSM, unit tests, a real `os/exec`-based executor with proper timeout/cancellation handling, a CLI, a goroutine/channel worker pool, and a `cmd/server` binary with real graceful shutdown (`SIGINT`/`SIGTERM` → drain workers → exit) — verified live. Accepted, known gaps going into V2: `internal/scheduler` has no unit tests yet (would need `executor` extracted to an interface first), and "worker registration" is only a loose in-process stand-in until real network registration arrives in V2/V3.

## What you should walk away knowing

- **Interfaces & structural typing** — small interfaces, satisfied implicitly, no `implements` keyword; why the `Repository` interface stays swappable for Postgres later.
- **Error handling idioms** — `%w` wrapping vs `%v`, `errors.Is`/sentinel errors, why the *thing you wrap* has to be the thing a caller will later check for.
- **Concurrency primitives** — goroutines, channels as a queue, the worker-pool pattern (shared channel, N consumers, no manual assignment logic needed), `sync.Mutex`/`RWMutex`, `sync.WaitGroup`, and the classic footguns: copying a struct that holds a mutex, closures capturing loop variables, `:=` shadowing inside nested blocks instead of reusing an outer variable.
- **`context.Context`** — as the root-context idiom, `WithTimeout`/cancellation propagation, and the non-obvious part: a killed subprocess's own error doesn't itself carry `context.DeadlineExceeded` — you have to check the context directly and deliberately wrap it if a caller needs to detect it later.
- **Manual dependency injection** — no framework; wiring happens by hand in `main()`, and every layer only depends on the interface it needs, not concrete types.
- **Table-driven testing** — plus the discipline of deriving expected values from the *spec*, not from the implementation being tested (a mistake that surfaced and got caught more than once here).
- **Package design** — organizing by domain (`job`, `scheduler`, `executor`) instead of generic `utils`/`helpers`.

Once V2-V4 land, add: protobuf/gRPC service boundaries, bidirectional streaming, Postgres persistence, Docker Compose, CI/CD.

## Is this showcase-worthy (Twitter/LinkedIn/GitHub) right now?

**Honestly: not yet, but it's on a real trajectory to be.** As it stands (V1 only — single process, in-memory, no networking), it demonstrates solid Go fundamentals but isn't architecturally distinct from a lot of portfolio repos out there. It's a genuinely good *checkpoint*, not the finish line to post about.

What would actually make it worth a public post:
- **V2 (gRPC + protobuf)** is the big one — real service boundaries and a typed contract are what separate "I built a queue" from "I built a distributed system," and it's the most resume-visible piece per plan.md's own framing.
- **V4 (Docker Compose + CI/CD + Postgres + integration tests)** is what makes it look production-grade rather than a weekend exercise.
- A proper root `README.md` with an architecture diagram, setup instructions, and example CLI usage — first impression matters as much as the code.
- More granular git history. Right now there's essentially one commit covering all of this work — a recruiter or curious viewer skimming commit history won't see the actual iterative build process, which is one of the strongest signals you actually understand what you built (vs. it being AI-authored wholesale).

**The genuine differentiator you'll have, once you do post it:** you can actually explain every design decision in an interview — why the repository is an interface, why the FSM is a map instead of scattered `if`s, why `%w` matters, why the worker pool needs no manual assignment logic. That's real and rare; most portfolio projects (AI-assisted or not) can't survive that kind of follow-up questioning. Worth saying explicitly in a post/README, not just implying it.

**Recommendation:** hold off on publicizing until at least V2 is done (ideally V4), keep committing incrementally from here on, and write the README last, once there's a finished story to tell.
