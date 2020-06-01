package xdg

import (
	"os"
	"path/filepath"
)

func getXDGDirectory(envVar string, fallback string) string {
	fromEnv := os.Getenv(envVar)
	if fromEnv != "" && filepath.IsAbs(fromEnv) {
		return fromEnv
	}
	return os.ExpandEnv(fallback)
}

// GetConfigHome returns the path of $XDG_DATA_HOME directory
func GetConfigHome() string {
	return getXDGDirectory("XDG_CONFIG_HOME", "$HOME/.config")
}
