// @feature:watch File modification time tracking for incremental builds.
package watch

import (
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

// MtimeStore tracks file modification times so that successive calls to
// Update can report which files changed or were removed.
type MtimeStore struct {
	Entries map[string]time.Time // key = RelPath
}

func NewMtimeStore() *MtimeStore {
	return &MtimeStore{
		Entries: make(map[string]time.Time),
	}
}

// Update walks inputDir and compares each file's modification time against the
// previously stored value. It returns the relative paths of changed (new or
// modified) and removed files. Dotfiles and _hidden_ prefixed paths are
// skipped, matching obsidian.Scan behaviour.
func (s *MtimeStore) Update(inputDir string) (changed []string, removed []string, err error) {
	seen := make(map[string]struct{})

	err = filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}

		relPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return err
		}

		if path == inputDir {
			return nil
		}

		name := d.Name()
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_hidden_") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		seen[relPath] = struct{}{}
		modTime := info.ModTime()

		prev, tracked := s.Entries[relPath]
		if !tracked || !modTime.Equal(prev) {
			changed = append(changed, relPath)
		}
		s.Entries[relPath] = modTime

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	for relPath := range s.Entries {
		if _, ok := seen[relPath]; !ok {
			removed = append(removed, relPath)
			delete(s.Entries, relPath)
		}
	}

	return changed, removed, nil
}
