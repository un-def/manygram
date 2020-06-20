package cli

import (
	"strings"

	"github.com/un-def/manygram/internal/profile"
	"github.com/un-def/manygram/internal/tg"
)

func init() {
	configCommand.AddCommand(
		"check", "Check the config", "Check the config.",
		new(configCheckCmd),
	)
}

type configCheckCmd struct{}

func (c *configCheckCmd) Execute(args []string) error {
	conf, err := readConfig()
	if err != nil {
		return err
	}
	printMessage("Config %s found. Checking.", getConfigPath())
	telegram, err := tg.Executable(conf.ExecPath)
	if err != nil {
		return newError("Check error: `exec-path`", err)
	}
	execMsg := []string{telegram.Path}
	if telegram.FullPath != telegram.Path {
		execMsg = append(execMsg, telegram.FullPath)
	}
	if telegram.RealPath != telegram.FullPath {
		execMsg = append(execMsg, telegram.RealPath)
	}
	printMessage("Telegram Desktop executable: %s", strings.Join(execMsg, " -> "))
	profileDirExist, err := profile.IsProfileDirExist(conf.ProfileDir)
	if err != nil {
		return newError("Check error: `profile-dir`", err)
	}
	profile.IsProfileDirExist(conf.ProfileDir)
	printMessage("Profile directory: %s", conf.ProfileDir)
	if !profileDirExist {
		printMessage("Profile directory does not exist.")
	}
	printMessage("OK. Check passed.")
	return nil
}
