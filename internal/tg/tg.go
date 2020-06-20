package tg

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func wrapError(err error) error {
	return fmt.Errorf("check executable: %w", err)
}

// DefaultPath is the default path/name of Telegram Desktop executable
const DefaultPath = "telegram-desktop"

// TelegramDesktop represents Telegram Desktop executable
type TelegramDesktop struct {
	Path     string
	FullPath string
	RealPath string
}

// Executable returns TelegramDesktop struct or error if executable not found
func Executable(path string) (*TelegramDesktop, error) {
	fullPath, err := exec.LookPath(path)
	if err != nil {
		return nil, wrapError(err)
	}
	realPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		return nil, wrapError(err)
	}
	return &TelegramDesktop{path, fullPath, realPath}, nil
}

// Run executes telegram-desktop executable
func (tg *TelegramDesktop) Run(profilePath string, extraArgs []string, wait bool) error {
	args := append([]string{tg.Path, "-many", "-workdir", profilePath}, extraArgs...)
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
