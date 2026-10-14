---
name: go-concurrent-worker-pool
description: Design and review bounded goroutine worker pools in Go — fan-out/fan-in over channels, graceful shutdown via a quit channel and sync.WaitGroup, and sizing workers to runtime.NumCPU or a container's CPU share. Use when building concurrent background processing, parallelizing CPU-bound work per frame or per batch, capping concurrency against a database or downstream API, or reviewing goroutine/channel code for leaks and unbounded concurrency.
---

# Go concurrent worker pool

A worker pool runs a **fixed** number of long-lived goroutines that pull work from
channels, instead of spawning one goroutine per task. This caps concurrency,
protects downstream systems, and reuses goroutines across many units of work.

## When to reach for it

- CPU-bound work you want spread across cores (per-frame entity updates, image
  processing, parsing a large batch).
- I/O against a resource with limited capacity (a DB that handles 100 QPS, a rate-limited
  API). A goroutine-per-request design will exhaust connections, file descriptors, or the
  remote service. A bounded pool applies **backpressure** instead.

## The pattern (persistent workers, signalled per batch)

From the game engine's parallel entity manager (`enemy/parallel.go` in git history):

```go
// One long-lived goroutine per shard. It blocks until signalled, does its slice
// of work, then reports done. A closed quit channel unwinds every worker.
func (em *Manager) worker(id int) {
    defer em.wg.Done()
    for {
        select {
        case <-em.quit:              // graceful shutdown
            return
        case <-em.workSignal[id]:    // "process this batch"
            em.shards[id].Update()
            em.done <- struct{}{}    // fan back in
        }
    }
}

// Update: fan out to all workers, then wait for all of them.
func (em *Manager) Update() {
    for i := range em.workSignal {
        em.workSignal[i] <- struct{}{}
    }
    for range em.workSignal {
        <-em.done
    }
}
```

For the common **jobs-channel** variant, prefer a single buffered `jobs` channel that
all workers range over, plus a `results` channel:

```go
jobs := make(chan Job, queueDepth)   // buffer sets the backpressure point
results := make(chan Result)
for i := 0; i < workers; i++ {
    go func() {
        for j := range jobs {        // exits when jobs is closed
            results <- process(j)
        }
    }()
}
```

## Sizing the pool

- CPU-bound: `runtime.NumCPU()` (leave one core for the main loop if it is hot).
- Container-aware: a container's CPU *limit* is not `NumCPU()`. Detect the environment or
  read the cgroup quota and cap accordingly (the game drops to 1 worker under Docker).
- I/O-bound: size to the downstream's safe concurrency, not the core count.

## Checklist

- [ ] One owner is responsible for closing each channel; close `jobs` when no more work
      will be sent, never from a worker.
- [ ] A `context.Context` or `quit` channel can cancel every worker; workers select on it.
- [ ] `sync.WaitGroup` (or draining `done`) lets shutdown block until workers exit.
- [ ] The jobs channel is buffered so producers block (backpressure) instead of growing memory.
- [ ] Shared state touched by workers is guarded by a mutex or split so shards don't overlap.

## Pitfalls

- **Goroutine leaks**: a worker blocked on a send/receive with no reader/quit path never
  exits. Always give it a cancel path.
- **Unbounded goroutines**: `go handle(req)` per request is not a pool; under load it melts.
- **Sending on a closed channel** panics. The producer closes; consumers only read.

## Backend relevance

This is the core concurrency primitive behind job queues, webhook dispatchers, batch
importers, and any service that must stay within a downstream's capacity. Demonstrating a
correct pool — with cancellation and backpressure — is a direct signal of backend maturity.

## References

- [Worker pool pattern in Go](https://github.com/romangurevitch/ConcurrencyWorkshop/blob/main/internal/pattern/workerpool/README.md)
- [Go concurrency: goroutines and channels](https://oneuptime.com/blog/post/2026-02-01-go-concurrency-goroutines-channels/view)
