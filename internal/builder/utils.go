package builder

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Init checks if the input directory (vault) exists.
// If not, it creates the directory and a default "Home.md" welcome note.
func Init(log *slog.Logger) {
	_, err := os.Stat(InputDir)

	if err == nil {
		log.Error("Vault directory already exists")
		return
	}

	if !os.IsNotExist(err) {
		log.Error("Couldn't read information about directory", "error", err)
		return
	}

	err = os.Mkdir(InputDir, 0755)
	if err != nil {
		log.Error("Couldn't create folder", "error", err)
		return
	}

	log.Info("Created vault directory")

	// Create a welcome note to get the user started
	welcomeText := "# Welcome to Kiln\n\nThis is your new vault. Run `kiln generate` to build it!"
	err = os.WriteFile(filepath.Join(InputDir, "Home.md"), []byte(welcomeText), 0644)
	if err != nil {
		log.Error("Couldn't create welcome note", "error", err)
		return
	}

	log.Info("Initialization complete")
}

// CleanOutputDir removes the entire output directory to ensure a clean build.
// This prevents stale files from persisting in the generated site.
func CleanOutputDir(log *slog.Logger) {
	err := os.RemoveAll(OutputDir)
	if err != nil {
		log.Error("Couldn't remove output directory", "error", err)
	} else {
		log.Info("Cleaned directory", "path", OutputDir)
	}
}

// isImageExt checks if the given file extension corresponds to a supported image format.
func isImageExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
		return true
	default:
		return false
	}
}

// isAllowedExt checks if the given extension corresponds to an allowed format into default mode generation
func isAllowedExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".pdf", ".ico":
		return true
	default:
		return false
	}
}
