// @feature:watch Tests for filesystem watcher.
package watch

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatchDetectsFileChange(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.md"), "initial")

	var called atomic.Int32
	done := make(chan struct{}, 1)

	w := &Watcher{
		InputDir: dir,
		Debounce: 50 * time.Millisecond,
		OnRebuild: func() error {
			called.Add(1)
			select {
			case done <- struct{}{}:
			default:
			}
			return nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- w.Watch(ctx) }()

	// Give the watcher time to start.
	time.Sleep(100 * time.Millisecond)

	writeFile(t, filepath.Join(dir, "a.md"), "modified")

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("callback did not fire after file change")
	}

	if n := called.Load(); n < 1 {
		t.Errorf("callback called %d times, want >= 1", n)
	}

	cancel()
	if err := <-errCh; err != nil {
		t.Errorf("Watch returned error: %v", err)
	}
}

func TestWatchDebounce(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.md"), "v0")

	var called atomic.Int32
	done := make(chan struct{}, 1)

	w := &Watcher{
		InputDir: dir,
		Debounce: 200 * time.Millisecond,
		OnRebuild: func() error {
			called.Add(1)
			select {
			case done <- struct{}{}:
			default:
			}
			return nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- w.Watch(ctx) }()

	time.Sleep(100 * time.Millisecond)

	// 5 rapid writes — should debounce into a single callback.
	for i := range 5 {
		writeFile(t, filepath.Join(dir, "a.md"), string(rune('a'+i)))
		time.Sleep(20 * time.Millisecond)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("callback did not fire")
	}

	// Wait a bit longer to confirm no extra callbacks fire.
	time.Sleep(300 * time.Millisecond)

	if n := called.Load(); n != 1 {
		t.Errorf("callback called %d times, want 1", n)
	}

	cancel()
	if err := <-errCh; err != nil {
		t.Errorf("Watch returned error: %v", err)
	}
}

func TestWatchContextCancellation(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.md"), "data")

	w := &Watcher{
		InputDir: dir,
		Debounce: 50 * time.Millisecond,
		OnRebuild: func() error {
			t.Error("callback should not fire")
			return nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() { errCh <- w.Watch(ctx) }()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Watch returned error: %v, want nil", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Watch did not return after context cancellation")
	}
}

func TestWatchNewSubdirectory(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.md"), "root")

	var called atomic.Int32
	done := make(chan struct{}, 1)

	w := &Watcher{
		InputDir: dir,
		Debounce: 50 * time.Millisecond,
		OnRebuild: func() error {
			called.Add(1)
			select {
			case done <- struct{}{}:
			default:
			}
			return nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- w.Watch(ctx) }()

	time.Sleep(100 * time.Millisecond)

	// Create a new subdirectory and file — watcher should detect it.
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)
	writeFile(t, filepath.Join(sub, "b.md"), "nested")

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("callback did not fire for new subdirectory file")
	}

	cancel()
	if err := <-errCh; err != nil {
		t.Errorf("Watch returned error: %v", err)
	}
}
