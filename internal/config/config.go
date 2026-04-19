package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Process struct {
	Cmd string            `toml:"cmd"`
	Env map[string]string `toml:"env"`
}

type Layout struct {
	Processes []string          `toml:"processes"`
	Env       map[string]string `toml:"env"`
}

type Config struct {
	Processes map[string]Process `toml:"processes"`
	Layouts   map[string]Layout  `toml:"layouts"`
}

func Load() (*Config, error) {
	path, err := findConfig()
	if err != nil {
		return nil, err
	}
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &cfg, nil
}

func findConfig() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		p := filepath.Join(dir, "proc.toml")
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("proc.toml not found")
}

// Resolve returns the ordered list of process names for a given set of args.
// Args can be a single layout name or a list of process names.
func (c *Config) Resolve(args []string) ([]string, map[string]string, error) {
	if len(args) == 1 {
		if layout, ok := c.Layouts[args[0]]; ok {
			return layout.Processes, layout.Env, nil
		}
	}
	for _, name := range args {
		if _, ok := c.Processes[name]; !ok {
			return nil, nil, fmt.Errorf("unknown process or layout: %q", name)
		}
	}
	return args, nil, nil
}
