# GoFleet - Project Context & Design Specification

## Overview

This document captures the complete context of the project being built so that another AI agent (or developer) can immediately understand the goals, philosophy, architecture, and implementation roadmap without requiring the original conversation.

---

# Primary Goal

Build a **real-world backend project in Go** that teaches:

- Go language idioms
- Concurrency
- Goroutines
- Channels
- Context cancellation
- Interfaces
- Error handling
- Package organization
- Dependency injection
- gRPC
- Streaming RPCs
- System Design
- Production engineering practices
- Graceful shutdown
- Testing
- Docker
- CI/CD

The project should not be a simple CRUD application.

It should resemble a real production backend system that could reasonably exist inside a tech company.

---

# Project Name

GoFleet

Alternative names are acceptable, but the architecture should remain the same.

---

# High-Level Idea

A distributed job execution platform.

Users submit jobs.

A scheduler assigns jobs to available workers.

Workers execute jobs.

Workers stream logs back.

Users can watch logs live.

Think of it as a simplified combination of:

- GitHub Actions Runner
- Celery
- Kubernetes Jobs
- Jenkins Agent
- Distributed Worker Queue

---

# Why This Project?

This project naturally exercises nearly every important Go concept.

Instead of forcing concurrency into a CRUD app, concurrency becomes the foundation of the architecture.

It also creates an excellent portfolio project for backend engineering interviews.

---

# Communication Choice

The project should primarily use **gRPC**.

Raw TCP sockets are **not** the primary transport.

Reasoning:

- strongly typed contracts
- protobuf
- unary RPC
- server streaming
- client streaming
- bidirectional streaming
- interceptors
- deadlines
- production-grade service communication

Raw sockets may later be implemented as an educational module to understand what gRPC abstracts away.

---

# System Architecture

```
CLI
 │
 │ gRPC
 ▼

Control Plane

 ├── Job API
 ├── Scheduler
 ├── Storage
 ├── Log Broker

        │
        │ gRPC Streaming
        ▼

Worker Agent
```

There are three executables:

```
server
worker
cli
```

Do NOT split into microservices initially.

Keep deployment simple.

---

# User Flow

1. User submits a job.
2. Scheduler stores job.
3. Worker connects.
4. Scheduler assigns job.
5. Worker executes.
6. Worker streams logs.
7. Scheduler stores state.
8. User watches logs.
9. Job finishes.

---

# Example Job

```json
{
  "command": "ping google.com",
  "timeout_seconds": 30,
  "priority": "high"
}
```

---

# Example CLI

Submit

```bash
gofleet submit --cmd "echo hello"
```

Status

```bash
gofleet status --job-id 123
```

Logs

```bash
gofleet logs --job-id 123 --follow
```

Cancel

```bash
gofleet cancel --job-id 123
```

---

# Recommended Repository Structure

```
gofleet/

cmd/
    server/
    worker/
    cli/

internal/

    job/
        model.go
        service.go
        repository.go

    scheduler/
        scheduler.go
        assignment.go
        heartbeat.go

    worker/
        executor.go
        runner.go
        stream.go

    transport/
        grpc/
            job_server.go
            worker_server.go
            interceptors.go

    storage/
        postgres.go
        migrations/

    logs/
        broker.go
        subscriber.go

    config/
        config.go

api/
    proto/
        gofleet.proto

gen/

pkg/
    errors/

scripts/

docker-compose.yml

Makefile

go.mod
```

---

# Why This Folder Structure?

Organize around business domains.

Avoid generic folders like

```
utils
helpers
common
misc
```

The package layout should communicate business intent.

---

# Transport Layer

Use gRPC.

Implement:

## JobService

- SubmitJob
- GetJob
- CancelJob
- WatchLogs

## WorkerService

- Connect

WorkerService should use bidirectional streaming.

---

# Database

Initially:

In-memory repository.

Later:

Postgres (preferred)

SQLite acceptable during development.

---

# Core Concepts to Learn

## Goroutines

Examples:

- scheduler loop
- heartbeat sender
- worker listener
- log broadcaster

Understand:

- lightweight concurrency
- scheduling
- synchronization

---

## Channels

Use channels for:

- job assignment
- log fan-out
- cancellation
- scheduler events

Do NOT abuse channels.

Shared mutable state should generally use mutexes.

---

## Context

Every long-running function should accept:

```go
context.Context
```

Examples:

```go
Run(ctx)

Execute(ctx)

CreateJob(ctx)

GetJob(ctx)
```

Cancellation should propagate through the system.

---

## Interfaces

Prefer small interfaces.

Example:

```go
type JobRepository interface {

    Create(...)

    GetByID(...)

    UpdateStatus(...)
}
```

Avoid giant "God interfaces."

---

## Error Handling

Wrap errors.

Example:

```go
return fmt.Errorf("create job: %w", err)
```

Use:

```
errors.Is()
errors.As()
```

Avoid exception-style thinking.

---

## Graceful Shutdown

Support:

- SIGINT
- SIGTERM
- active job cleanup
- worker disconnect
- server shutdown

---

# Job Lifecycle

Allowed transitions:

```
PENDING

↓

RUNNING

↓

SUCCESS
FAILED
TIMEOUT
CANCELLED
```

Failed jobs may return to:

```
FAILED

↓

PENDING
```

if retries remain.

Never allow:

```
SUCCESS
ok 
↓

RUNNING
```

or other invalid transitions.

Treat job state as a finite-state machine.

---

# Version 1 Features

Must include:

- CLI
- Job submission
- Scheduler
- Worker
- In-memory repository
- Job execution
- Job status
- Unit tests

No streaming required yet.

---

# Version 2

Introduce gRPC.

Add:

- protobuf
- JobService
- WorkerService
- generated code
- CLI client

---

# Version 3

Introduce streaming.

Features:

- live logs
- heartbeat
- worker registration
- job assignment
- timeout
- cancellation

---

# Version 4

Production improvements.

Include:

- Postgres
- retries
- graceful shutdown
- Docker Compose
- integration tests
- GitHub Actions

---

# Stretch Goals

After the project is stable:

- priority queue
- multiple workers
- worker tags
- Prometheus metrics
- OpenTelemetry tracing
- authentication
- rate limiting
- Docker sandbox execution
- web dashboard

---

# Design Principles

## 1.

Business-oriented packages.

---

## 2.

Keep interfaces close to consumers.

---

## 3.

Context everywhere.

---

## 4.

No global state.

Use dependency injection.

---

## 5.

Explicit state transitions.

---

## 6.

Cancellation should actually terminate running work.

---

## 7.

Write tests throughout development instead of only at the end.

---

## 8.

Avoid premature abstraction.

Refactor only after duplication becomes obvious.

---

# Why gRPC Instead of Raw Sockets?

Primary learning objectives include:

- protobuf
- streaming
- interceptors
- typed contracts
- service boundaries

Raw sockets teach protocol design but are lower priority.

After the gRPC version is complete, implement a second transport using raw TCP solely for educational purposes.

Example protocol:

```
WORKER_HELLO worker-1

HEARTBEAT worker-1

RUN job123 echo hello

LOG job123 hello

DONE job123 SUCCESS
```

This demonstrates what gRPC abstracts internally.

---

# Resume Value

This project should ultimately be described as:

> Built GoFleet, a distributed job execution platform in Go using gRPC, Protocol Buffers, bidirectional streaming, worker heartbeats, context-based cancellation, live log streaming, retries, graceful shutdown, and persistent job state management.

Potential resume bullets:

- Designed scheduler-worker architecture using gRPC bidirectional streams.
- Implemented real-time log streaming.
- Added worker heartbeats and failure detection.
- Implemented retries and cancellation using Go contexts.
- Built Docker Compose deployment.
- Added CI/CD with GitHub Actions.
- Wrote unit and integration tests.

---

# Learning Philosophy

The goal is **not** merely to finish a project.

The goal is to deeply understand:

- Go idioms
- Production backend engineering
- Concurrency
- Networking
- Package design
- System design
- Error handling
- Maintainable architecture

Every architectural decision should favor learning long-term Go engineering practices rather than simply making the project work.

---

# Success Criteria

A completed Version 1 should support:

- submitting jobs
- executing jobs
- viewing status
- worker registration
- clean package organization
- unit tests

A completed Version 2 should additionally support:

- gRPC
- protobuf
- streaming
- live logs
- cancellation
- heartbeats

A completed production version should additionally support:

- persistence
- retries
- graceful shutdown
- Docker Compose
- CI/CD
- metrics
- integration tests

The final result should resemble a production-grade distributed backend system rather than a tutorial project.


# IMPORTANT
 I want to use this project to get fluent in Go atleast in the sense that basic api integrations and functional setups come normally to me, so don;t write a single piece of code, you role is to guide me via snippets and diagrams, or citations to weblinks helping understand concept.