package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// Config holds application configuration paths.
type Config struct {
	DataDir string
	DBPath  string
}

// Load returns the application config with platform-appropriate paths.
func Load() (*Config, error) {
	dataDir, err := appDataDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	return &Config{
		DataDir: dataDir,
		DBPath:  filepath.Join(dataDir, "pulsereader.db"),
	}, nil
}

func appDataDir() (string, error) {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "PulseReader"), nil
	}

	// Linux / macOS: use XDG_DATA_HOME or ~/.local/share
	xdg := os.Getenv("XDG_DATA_HOME")
	if xdg == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		xdg = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(xdg, "pulsereader"), nil
}
