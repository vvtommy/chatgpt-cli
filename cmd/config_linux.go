package main

import (
	"os"
	"path/filepath"
)

var configFolder = func() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return os.Getenv("XDG_CONFIG_HOME")
	}
	return filepath.Join(os.Getenv("HOME"), ".config")
}()
