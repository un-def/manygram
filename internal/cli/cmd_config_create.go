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
	var dataDir string
	telegram, err := tg.Executable(tg.DefaultPath, nil)
	if err == nil {
		printMessage("Telegram Desktop executable: %s -> %s", telegram.Path, telegram.FullPath)
		conf.ExecPath = telegram.Path
		if telegram.IsSnap() {
			printMessage("Telegram Desktop seems installed via snap.")
			dataDir, err = tg.GetSnapDataHome()
			if err != nil {
				printMessage("Cannot find snap data directory, use fallback data location.")
				// snap has no access to dot files/dirs (e.g., ~/.local)
				dataDir = os.Getenv("HOME")
			}
		}
	} else {
		telegram, err = tg.Flatpak()
		if err == nil {
			printMessage("Telegram Desktop installed via Flatpak.")
			conf.ExecPath = telegram.Path
			conf.ExecArgs = telegram.Args
			dataDir, err = tg.GetFlatpakDataHome()
			if err != nil {
				printMessage("Cannot find Flatpak data directory, use fallback data location.")
				dataDir = xdg.GetDataHome()
			}
		} else {
			printMessage("Telegram Desktop executable not found.")
		}
	}
	if conf.ExecPath == "" {
		conf.ExecPath = tg.DefaultPath
	}
	if dataDir == "" {
		dataDir = xdg.GetDataHome()
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
