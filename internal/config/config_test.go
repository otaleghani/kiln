// @feature:cli
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kiln.yaml")
	content := `theme: dracula
font: merriweather
url: https://example.com
name: My Site
input: ./notes
output: ./dist
mode: custom
layout: simple
flat-urls: true
disable-toc: true
disable-local-graph: true
port: "3000"
host: 127.0.0.1
log: debug
accent-color: blue
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write kiln.yaml: %v", err)
	}

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
	if cfg.Font != "merriweather" {
		t.Errorf("Font = %q, want %q", cfg.Font, "merriweather")
	}
	if cfg.URL != "https://example.com" {
		t.Errorf("URL = %q, want %q", cfg.URL, "https://example.com")
	}
	if cfg.Name != "My Site" {
		t.Errorf("Name = %q, want %q", cfg.Name, "My Site")
	}
	if cfg.Input != "./notes" {
		t.Errorf("Input = %q, want %q", cfg.Input, "./notes")
	}
	if cfg.Output != "./dist" {
		t.Errorf("Output = %q, want %q", cfg.Output, "./dist")
	}
	if cfg.Mode != "custom" {
		t.Errorf("Mode = %q, want %q", cfg.Mode, "custom")
	}
	if cfg.Layout != "simple" {
		t.Errorf("Layout = %q, want %q", cfg.Layout, "simple")
	}
	if !cfg.FlatURLs {
		t.Error("FlatURLs = false, want true")
	}
	if !cfg.DisableTOC {
		t.Error("DisableTOC = false, want true")
	}
	if !cfg.DisableLocalGraph {
		t.Error("DisableLocalGraph = false, want true")
	}
	if cfg.Port != "3000" {
		t.Errorf("Port = %q, want %q", cfg.Port, "3000")
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("Host = %q, want %q", cfg.Host, "127.0.0.1")
	}
	if cfg.Log != "debug" {
		t.Errorf("Log = %q, want %q", cfg.Log, "debug")
	}
	if cfg.AccentColor != "blue" {
		t.Errorf("AccentColor = %q, want %q", cfg.AccentColor, "blue")
	}
}

func TestLoad_PartialFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kiln.yaml")
	content := `theme: nord
url: https://notes.dev
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write kiln.yaml: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	if cfg.Theme != "nord" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "nord")
	}
	if cfg.URL != "https://notes.dev" {
		t.Errorf("URL = %q, want %q", cfg.URL, "https://notes.dev")
	}

	// All other fields should be zero values.
	if cfg.Font != "" {
		t.Errorf("Font = %q, want empty", cfg.Font)
	}
	if cfg.Name != "" {
		t.Errorf("Name = %q, want empty", cfg.Name)
	}
	if cfg.Input != "" {
		t.Errorf("Input = %q, want empty", cfg.Input)
	}
	if cfg.Output != "" {
		t.Errorf("Output = %q, want empty", cfg.Output)
	}
	if cfg.Mode != "" {
		t.Errorf("Mode = %q, want empty", cfg.Mode)
	}
	if cfg.Layout != "" {
		t.Errorf("Layout = %q, want empty", cfg.Layout)
	}
	if cfg.FlatURLs {
		t.Error("FlatURLs = true, want false")
	}
	if cfg.DisableTOC {
		t.Error("DisableTOC = true, want false")
	}
	if cfg.DisableLocalGraph {
		t.Error("DisableLocalGraph = true, want false")
	}
	if cfg.Port != "" {
		t.Errorf("Port = %q, want empty", cfg.Port)
	}
	if cfg.Host != "" {
		t.Errorf("Host = %q, want empty", cfg.Host)
	}
	if cfg.Log != "" {
		t.Errorf("Log = %q, want empty", cfg.Log)
	}
	if cfg.AccentColor != "" {
		t.Errorf("AccentColor = %q, want empty", cfg.AccentColor)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	cfg, err := Load("/nonexistent/path/kiln.yaml")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if cfg != nil {
		t.Fatalf("expected nil config, got %+v", cfg)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kiln.yaml")
	if err := os.WriteFile(path, []byte(":::bad yaml[[["), 0o644); err != nil {
		t.Fatalf("write kiln.yaml: %v", err)
	}

	cfg, err := Load(path)
	if err == nil {
		t.Fatal("expected non-nil error for invalid YAML")
	}
	if cfg != nil {
		t.Errorf("expected nil config on error, got %+v", cfg)
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kiln.yaml")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatalf("write kiln.yaml: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config for empty file")
	}

	// All fields should be zero values.
	if cfg.Theme != "" {
		t.Errorf("Theme = %q, want empty", cfg.Theme)
	}
	if cfg.FlatURLs {
		t.Error("FlatURLs = true, want false")
	}
	if cfg.DisableTOC {
		t.Error("DisableTOC = true, want false")
	}
	if cfg.DisableLocalGraph {
		t.Error("DisableLocalGraph = true, want false")
	}
	if cfg.AccentColor != "" {
		t.Errorf("AccentColor = %q, want empty", cfg.AccentColor)
	}
}

func TestFindFile_Exists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kiln.yaml")
	if err := os.WriteFile(path, []byte("theme: nord\n"), 0o644); err != nil {
		t.Fatalf("write kiln.yaml: %v", err)
	}

	got, err := FindFile(dir)
	if err != nil {
		t.Fatalf("FindFile: %v", err)
	}
	if got != path {
		t.Errorf("FindFile = %q, want %q", got, path)
	}
}

func TestFindFile_NotFound(t *testing.T) {
	dir := t.TempDir()

	got, err := FindFile(dir)
	if err != nil {
		t.Fatalf("FindFile: %v", err)
	}
	if got != "" {
		t.Errorf("FindFile = %q, want empty", got)
	}
}

func TestValueOr_OnlyOverridesEmpty(t *testing.T) {
	cfg := Config{Theme: "dracula", Font: ""}

	if got := cfg.ValueOr("theme", "default"); got != "dracula" {
		t.Errorf("ValueOr(theme) = %q, want %q", got, "dracula")
	}
	if got := cfg.ValueOr("font", "inter"); got != "inter" {
		t.Errorf("ValueOr(font) = %q, want %q", got, "inter")
	}

	// accent-color: config value wins when set.
	cfgColor := Config{AccentColor: "red"}
	if got := cfgColor.ValueOr("accent-color", ""); got != "red" {
		t.Errorf("ValueOr(accent-color) = %q, want %q", got, "red")
	}
	// host: config value wins when set.
	cfgHost := Config{Host: "0.0.0.0"}
	if got := cfgHost.ValueOr("host", "127.0.0.1"); got != "0.0.0.0" {
		t.Errorf("ValueOr(host) = %q, want %q", got, "0.0.0.0")
	}
	// host: fallback wins when empty.
	cfgNoHost := Config{}
	if got := cfgNoHost.ValueOr("host", "127.0.0.1"); got != "127.0.0.1" {
		t.Errorf("ValueOr(host) = %q, want %q", got, "127.0.0.1")
	}
}
