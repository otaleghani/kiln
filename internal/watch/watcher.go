// @feature:watch Filesystem watcher using fsnotify with debounced rebuild trigger.
package watch

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const DefaultDebounce = 300 * time.Millisecond

// RebuildFunc is the callback invoked when file changes are detected.
type RebuildFunc func() error

// Watcher monitors InputDir for filesystem changes and calls OnRebuild after
// a debounce period of inactivity.
type Watcher struct {
	InputDir  string
	Debounce  time.Duration
	Log       *slog.Logger
	OnRebuild RebuildFunc
}

// Watch creates an fsnotify watcher, recursively adds all directories under
// InputDir (skipping dotfiles and _hidden_ prefixed paths), and runs an event
// loop until ctx is cancelled. Write, Create, Remove, and Rename events reset
// a debounce timer; when the timer fires, OnRebuild is called. Newly created
// directories are automatically added to the watcher.
func (w *Watcher) Watch(ctx context.Context) error {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fsw.Close()

	if err := w.addDirsRecursive(fsw, w.InputDir); err != nil {
		return err
	}

	debounce := w.Debounce
	if debounce == 0 {
		debounce = DefaultDebounce
	}

	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()

	log := w.Log
	if log == nil {
		log = slog.Default()
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-fsw.Events:
			if !ok {
				return nil
			}

			if !isRelevantOp(event.Op) {
				continue
			}

			if shouldSkipPath(event.Name) {
				continue
			}

			if event.Op&fsnotify.Create != 0 {
				w.tryAddDir(fsw, event.Name, log)
			}

			timer.Reset(debounce)

		case watchErr, ok := <-fsw.Errors:
			if !ok {
				return nil
			}
			log.Warn("fsnotify error", "err", watchErr)

		case <-timer.C:
			log.Info("rebuilding after file change")
			if rebuildErr := w.OnRebuild(); rebuildErr != nil {
				log.Error("rebuild failed", "err", rebuildErr)
			}
		}
	}
}

func isRelevantOp(op fsnotify.Op) bool {
	return op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0
}

func shouldSkipPath(path string) bool {
	for _, segment := range strings.Split(filepath.ToSlash(path), "/") {
		if strings.HasPrefix(segment, ".") || strings.HasPrefix(segment, "_hidden_") {
			return true
		}
	}
	return false
}

func (w *Watcher) addDirsRecursive(fsw *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		if path != root && (strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_hidden_")) {
			return filepath.SkipDir
		}
		return fsw.Add(path)
	})
}

func (w *Watcher) tryAddDir(fsw *fsnotify.Watcher, path string, log *slog.Logger) {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return
	}
	if addErr := fsw.Add(path); addErr != nil {
		log.Warn("failed to watch new directory", "path", path, "err", addErr)
	}
}
