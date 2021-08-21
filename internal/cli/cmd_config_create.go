package cli

import (
	"os"

	"github.com/un-def/manygram/internal/config"
	"github.com/un-def/manygram/internal/tg"
	"github.com/un-def/manygram/internal/util"
	"github.com/un-def/manygram/internal/xdg"
)

func init() {
	configCommand.AddCommand(
		"create", "Create the config", "Create the config.",
		new(configCreateCmd),
	)
}

type configCreateCmd struct {
	Force bool `short:"f" long:"force" description:"Rewrite the existing config"`
}

func (c *configCreateCmd) Execute(args []string) error {
	configPath := getConfigPath()
	exist, err := util.Exist(configPath)
	if err != nil {
		return err
	}
	if !exist {
		printMessage("Config %s is not found. Creating a new one.", configPath)
	} else if c.Force {
		printMessage("Config %s has been found. Recreating.", configPath)
	} else {
		return newError("Config %s already exists.", configPath)
	}
	conf := config.New(configPath)
	execPath := tg.DefaultPath
	conf.ExecPath = execPath
	dataDir := xdg.GetDataHome()
	telegram, err := tg.Executable(execPath, nil)
	if err != nil {
		printMessage("Telegram Desktop executable not found.")
	} else {
		printMessage("Telegram Desktop executable: %s -> %s", telegram.Path, telegram.FullPath)
	}
	if telegram != nil && telegram.IsSnap() {
		printMessage("Telegram Desktop seems installed via snap.")
		if dataDir, err = tg.GetSnapDataHome(); err != nil {
			printMessage("Cannot find snap data directory, use fallback data location.")
			dataDir = os.Getenv("HOME")
		}
	}
	profileDir := getDefaultProfileDir(dataDir)
	printMessage("Profile directory: %s", profileDir)
	conf.ProfileDir = profileDir
	if err = conf.Write(); err != nil {
		return newError("Failed to write config.", err)
	}
	printMessage("Done.")
	return nil
}
