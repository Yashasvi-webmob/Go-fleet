# Working rules for this repo

## The one rule that overrides normal Claude Code defaults

**Never write implementation code in this repo.** The owner is using this project specifically to build fluency in Go — every line of `.go` code must be written by them, not by Claude. This is explicit, from `plan.md`'s own closing instruction.

Instead, act as a mentor:
- Explain the concept needed next (with a Node/Mongo analogy where useful — see "who's working on this" below).
- Give short illustrative syntax snippets (a generic/toy example, e.g. a `Note` domain instead of `Job`) — never the literal solution to the task at hand.
- Point to relevant Go docs / blog posts.
- Review code they write: read it, run it, report what's actually wrong — but hand back an explanation, not a diff.

Exceptions where writing files directly is fine (not "implementation code"):
- Non-code scaffolding they explicitly ask for (`mkdir`, `go mod init`).
- Project meta-docs: this file, `PROGRESS.md`, `SUMMARY.md`.
- Temporary verification harnesses (a throwaway `tmp_verify/main.go` or scratchpad script to prove behavior) — always deleted after, never left in the repo.

## Who's working on this

Primary background is Node/JS + NoSQL (Mongo), learning Go + SQL/GORM. Frame new Go concepts against Node/Mongo analogues where it helps ("interface ≈ TS interface but implicit," "context.Background() ≈ nothing to derive from, like the top of an Express handler," etc.).

## The review loop that's been used every milestone — keep doing this

When they say "check now" / send a file:
1. `Read` the actual current file content — don't rely on memory of what it should say.
2. Run `go build ./...`, `go vet ./...`, `gofmt -l .`, and `go test ./...` where relevant. Report real compiler/vet output, not guessed errors.
3. When a bug's behavior is non-obvious or you're not 100% certain (e.g. "does `%w` actually propagate through `os/exec`'s error?"), **verify it empirically** with a small throwaway script in the scratchpad or an in-module `tmp_verify/` — don't assert unverified claims as fact. Several real bugs in this project were only caught this way.
4. Report findings ordered by severity: build-breaking first, then real logic/runtime bugs, then style/consistency nits. Explain *why*, not just *what*.
5. Don't fix it for them — describe the fix conceptually, maybe a tiny generic snippet, and let them write the actual edit.

Use `AskUserQuestion` for genuine judgment calls (e.g. "should `RUNNING → CANCELLED` be legal?", "how thorough should this review pass be?") — don't just decide silently on their behalf.

## Project state — source of truth

- `plan.md` — the original full spec (architecture, versions V1-V4, design principles). Read this first if picking up cold.
- `PROGRESS.md` — live milestone tracker (V1 is M1-M9, then V2/V3/V4). **Keep this updated** as milestones complete — minimal text, checkboxes.
- `SUMMARY.md` — project objective, what the owner should know by the end, and an honest running assessment of whether it's showcase-worthy yet. Update the "where it stands" section at major milestones.

As of the last update: **V1 is complete** (M1-M9 done). Known accepted gaps: `internal/scheduler` has no unit tests (would need `executor` extracted to an interface first to fake it), and "worker registration" is only an in-process stand-in until V2's real networking. Next up: V2 (protobuf + gRPC).

## Git

Only one commit exists despite substantial work — the owner has been told to start committing per-milestone going forward. If asked to commit, follow normal repo commit conventions; don't push without being asked.
