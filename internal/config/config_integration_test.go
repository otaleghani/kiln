// @feature:cli Integration tests for config file loading
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigIntegration_FullPipeline(t *testing.T) {
	dir := t.TempDir()
	content := `theme: dracula
font: lato
url: https://example.com
accent-color: green
host: 127.0.0.1
`
	if err := os.WriteFile(filepath.Join(dir, "kiln.yaml"), []byte(content), 0o644); err != nil {
		t.Fatalf("write kiln.yaml: %v", err)
	}

	// FindFile should discover the config.
	path, err := FindFile(dir)
	if err != nil {
		t.Fatalf("FindFile: %v", err)
	}
	if path == "" {
		t.Fatal("FindFile returned empty path, want kiln.yaml")
	}

	// Load should parse all fields.
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.Theme != "dracula" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "dracula")
	}
	if cfg.Font != "lato" {
		t.Errorf("Font = %q, want %q", cfg.Font, "lato")
	}
	if cfg.URL != "https://example.com" {
		t.Errorf("URL = %q, want %q", cfg.URL, "https://example.com")
	}
	if cfg.AccentColor != "green" {
		t.Errorf("AccentColor = %q, want %q", cfg.AccentColor, "green")
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("Host = %q, want %q", cfg.Host, "127.0.0.1")
	}

	// ValueOr: config value wins over fallback.
	if got := cfg.ValueOr("theme", "default"); got != "dracula" {
		t.Errorf("ValueOr(theme) = %q, want %q", got, "dracula")
	}

	// ValueOr: host config value wins.
	if got := cfg.ValueOr("host", "0.0.0.0"); got != "127.0.0.1" {
		t.Errorf("ValueOr(host) = %q, want %q", got, "127.0.0.1")
	}

	// ValueOr: accent-color config value wins.
	if got := cfg.ValueOr("accent-color", ""); got != "green" {
		t.Errorf("ValueOr(accent-color) = %q, want %q", got, "green")
	}

	// ValueOr: fallback wins for unset field.
	if got := cfg.ValueOr("mode", "default"); got != "default" {
		t.Errorf("ValueOr(mode) = %q, want %q", got, "default")
	}
}

func TestConfigIntegration_BoolFields(t *testing.T) {
	dir := t.TempDir()
	content := `flat-urls: true
disable-toc: true
`
	if err := os.WriteFile(filepath.Join(dir, "kiln.yaml"), []byte(content), 0o644); err != nil {
		t.Fatalf("write kiln.yaml: %v", err)
	}

	path, err := FindFile(dir)
	if err != nil {
		t.Fatalf("FindFile: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	if !cfg.FlatURLs {
		t.Error("FlatURLs = false, want true")
	}
	if !cfg.DisableTOC {
		t.Error("DisableTOC = false, want true")
	}
	if cfg.DisableLocalGraph {
		t.Error("DisableLocalGraph = true, want false")
	}

	// BoolOr: config value wins.
	if got := cfg.BoolOr("flat-urls", false); !got {
		t.Error("BoolOr(flat-urls) = false, want true")
	}

	// BoolOr: fallback wins for unset field.
	if got := cfg.BoolOr("disable-local-graph", false); got {
		t.Error("BoolOr(disable-local-graph) = true, want false")
	}
}
