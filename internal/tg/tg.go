package tg

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// DefaultPath is the default path/name of Telegram Desktop executable
const DefaultPath = "telegram-desktop"

// TelegramDesktop represents Telegram Desktop executable
type TelegramDesktop struct {
	Path     string
	FullPath string
	RealPath string
	Args     []string
}

// Executable returns TelegramDesktop struct or error if executable not found
func Executable(path string, args []string) (*TelegramDesktop, error) {
	fullPath, err := exec.LookPath(path)
	if err != nil {
		return nil, err
	}
	realPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		return nil, err
	}
	return &TelegramDesktop{path, fullPath, realPath, args}, nil
}

// Run executes telegram-desktop executable
func (tg *TelegramDesktop) Run(profilePath string, extraArgs []string, wait bool) error {
	args := make([]string, len(tg.Args)+len(extraArgs)+3)
	copy(args, tg.Args)
	args = append(args, "-many", "-workdir", profilePath)
	args = append(args, extraArgs...)
	cmd := exec.Command(tg.Path, args...)
	if wait {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return cmd.Start()
}

// IsSnap returns true if executable seems installed with snap
func (tg *TelegramDesktop) IsSnap() bool {
	return path.Base(tg.RealPath) == "snap"
}

// GetSnapDataHome returns XDG_DATA_HOME of Telegram Desktop snap or error
func GetSnapDataHome() (string, error) {
	snapDataHome := os.ExpandEnv("$HOME/snap/telegram-desktop/current/.local/share")
	if _, err := os.Stat(snapDataHome); err != nil {
		return "", err
	}
	return snapDataHome, nil
}

const flatpakExecName = "flatpak"
const flatpakAppID = "org.telegram.desktop"

// Flatpak returns TelegramDesktop struct representing Telegram Desktop Flatpak app
// or error if Flatpak or app is not installed
func Flatpak() (*TelegramDesktop, error) {
	flatpakExecPath, err := exec.LookPath(flatpakExecName)
	if err != nil {
		return nil, errors.New("flatpak executable not found")
	}
	if err := exec.Command(flatpakExecPath, "--user", "info", flatpakAppID).Run(); err == nil {
		return Executable(flatpakExecName, []string{"run", "--user", flatpakAppID})
	}
	if err := exec.Command(flatpakExecPath, "info", flatpakAppID).Run(); err == nil {
		return Executable(flatpakExecName, []string{"run", flatpakAppID})
	}
	return nil, errors.New(flatpakAppID + " flatpak app not found")
}

// GetFlatpakHome returns XDG_DATA_HOME of Telegram Desktop Flatpak app or error
func GetFlatpakDataHome() (string, error) {
	flatpakDataHome := os.ExpandEnv("$HOME/.var/app/" + flatpakAppID + "/data")
	if _, err := os.Stat(flatpakDataHome); err != nil {
		return "", err
	}
	return flatpakDataHome, nil
}
