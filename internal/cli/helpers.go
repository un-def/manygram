package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/un-def/manygram/internal/xdg"
)

const profileDirName = "profiles"

func getConfigPath() string {
	return path.Join(xdg.GetConfigHome(), "manygram", "config.toml")
}

func getDefaultProfileDir(dataDir string) string {
	return path.Join(dataDir, "manygram", profileDirName)
}

func printMessage(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
	fmt.Fprint(os.Stdout, "\n")
}
