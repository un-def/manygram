package cli

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/un-def/manygram/internal/config"
	"github.com/un-def/manygram/internal/xdg"
)

const profileDirName = "profiles"

func getDefaultProfileDir(dataDir string) string {
	return path.Join(dataDir, "manygram", profileDirName)
}

func printMessage(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
	fmt.Fprint(os.Stdout, "\n")
}

func getConfigPath() string {
	return path.Join(xdg.GetConfigHome(), "manygram", "config.toml")
}

func readConfig() (*config.Config, error) {
	configPath := getConfigPath()
	conf, err := config.Read(configPath)
	if err == nil {
		return conf, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, newError(
			"Config %s not found. Run `manygram config` to create a new config.",
			configPath, err,
		)
	}
	return nil, newError("Failed to read config %s", configPath, err)
}
