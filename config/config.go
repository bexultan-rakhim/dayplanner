package config

import (
	"os"
	"path/filepath"
)

const appName = "dayplanner"

func Dir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, appName)
	}

	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".config", appName)
	}

	wd, err := os.Getwd()
	if err != nil {
		return appName 
	}
	return filepath.Join(wd, appName)
}

