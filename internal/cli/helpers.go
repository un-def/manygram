package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/un-def/manygram/internal/xdg"
)

var configPath string

func getConfigPath() string {
	if configPath != "" {
		return configPath
	}
	configHome := xdg.GetConfigHome()
	configPath = path.Join(configHome, "manygram", "config.toml")
	return configPath
}

func printMessage(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
	fmt.Fprint(os.Stdout, "\n")
}
