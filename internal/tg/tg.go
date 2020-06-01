package tg

import (
	"os"
	"os/exec"
)

// Run executes telegram-desktop binary
func Run(binPath string, profilePath string, extraArgs []string, wait bool) error {
	args := append([]string{binPath, "-many", "-workdir", profilePath}, extraArgs...)
	cmd := exec.Command(binPath, args...)
	if wait {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return cmd.Start()
}
