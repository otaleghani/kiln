# Go Concurrency

## Overview
Safe concurrency patterns for Go: goroutines, channels, sync primitives, and common pitfalls.

## Guidelines
1. Always use `sync.WaitGroup` or `errgroup.Group` to wait for goroutines.
2. Use `context.Context` for cancellation — pass it as the first parameter.
3. Prefer `sync.Mutex` over channels for protecting shared state.
4. Use `chan struct{}` for signaling, not `chan bool`.
5. Always handle the done/cancel case in select statements.
6. Use `sync.Once` for one-time initialization, not manual flags.
7. Buffer channels when the producer shouldn't block: `make(chan T, bufSize)`.
8. Close channels from the sender side only — never from the receiver.

## Patterns
```go
// errgroup for concurrent operations with error propagation
g, ctx := errgroup.WithContext(ctx)
for _, item := range items {
    item := item // capture loop variable
    g.Go(func() error {
        return process(ctx, item)
    })
}
if err := g.Wait(); err != nil {
    return fmt.Errorf("processing failed: %w", err)
}
```

## Anti-Patterns
- Don't launch goroutines without a way to wait for them or cancel them.
- Don't use `time.Sleep` for synchronization — use channels or sync primitives.
- Don't read and write a map concurrently without a mutex.
