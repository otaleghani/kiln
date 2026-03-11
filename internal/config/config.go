// @feature:cli Configuration file loading and type definitions
package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DefaultFilename is the conventional name for the config file.
const DefaultFilename = "kiln.yaml"

type Config struct {
	Theme             string `yaml:"theme"`
	Font              string `yaml:"font"`
	URL               string `yaml:"url"`
	Name              string `yaml:"name"`
	Input             string `yaml:"input"`
	Output            string `yaml:"output"`
	Mode              string `yaml:"mode"`
	Layout            string `yaml:"layout"`
	FlatURLs          bool   `yaml:"flat-urls"`
	DisableTOC        bool   `yaml:"disable-toc"`
	DisableLocalGraph bool   `yaml:"disable-local-graph"`
	DisableBacklinks  bool   `yaml:"disable-backlinks"`
	Port              string `yaml:"port"`
	Log               string `yaml:"log"`
	Lang              string `yaml:"lang"`
	AccentColor       string `yaml:"accent-color"`
}

// Load reads a kiln.yaml file from the given path.
// Returns (nil, nil) if the file does not exist (not an error).
// Returns (*Config, nil) on success.
// Returns (nil, error) on parse failure.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		if errors.Is(err, io.EOF) {
			return &cfg, nil
		}
		return nil, err
	}
	return &cfg, nil
}

// FindFile looks for kiln.yaml in the given directory.
// Returns the full path if found, or "" if not present.
// Returns an error only for filesystem errors (not for missing file).
func FindFile(dir string) (string, error) {
	path := filepath.Join(dir, DefaultFilename)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return path, nil
}

// ValueOr returns the config string field if non-empty, otherwise the fallback.
func (c *Config) ValueOr(field, fallback string) string {
	var val string
	switch field {
	case "theme":
		val = c.Theme
	case "font":
		val = c.Font
	case "url":
		val = c.URL
	case "name":
		val = c.Name
	case "input":
		val = c.Input
	case "output":
		val = c.Output
	case "mode":
		val = c.Mode
	case "layout":
		val = c.Layout
	case "port":
		val = c.Port
	case "log":
		val = c.Log
	case "lang":
		val = c.Lang
	case "accent-color":
		val = c.AccentColor
	}
	if val != "" {
		return val
	}
	return fallback
}

// BoolOr returns the config bool field if the config is non-nil, otherwise the fallback.
func (c *Config) BoolOr(field string, fallback bool) bool {
	switch field {
	case "flat-urls":
		return c.FlatURLs
	case "disable-toc":
		return c.DisableTOC
	case "disable-local-graph":
		return c.DisableLocalGraph
	case "disable-backlinks":
		return c.DisableBacklinks
	}
	return fallback
}
