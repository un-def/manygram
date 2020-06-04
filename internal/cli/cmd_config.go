package cli

import (
	"errors"
	"os"
	"strings"

	"github.com/un-def/manygram/internal/config"
	"github.com/un-def/manygram/internal/profile"
	"github.com/un-def/manygram/internal/tg"
	"github.com/un-def/manygram/internal/xdg"
)

func init() {
	parser.AddCommand(
		"config", "Show configuration",
		`Show and check a config. If the config does not exist, it will be generated.`,
		new(configCmd),
	)
}

type configCmd struct {
}

func (c *configCmd) Execute(args []string) error {
	configPath := getConfigPath()
	conf, err := config.Read(configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return newError("Failed to read config %s", configPath, err)
		}
		printMessage("Config %s not found. Creating a new one.", configPath)
		conf = config.New(configPath)
		conf.ExecPath = tg.DefaultPath
		dataDir := xdg.GetDataHome()
		telegram, err := tg.Executable(tg.DefaultPath)
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
	printMessage("Config %s found. Checking.", configPath)
	telegram, err := tg.Executable(conf.ExecPath)
	if err != nil {
		return newError(`Check error: "exec-path"`, err)
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
		return newError(`Check error: "profile-dir"`, err)
	}
	profile.IsProfileDirExist(conf.ProfileDir)
	printMessage("Profile directory: %s", conf.ProfileDir)
	if !profileDirExist {
		printMessage("Profile directory does not exist.")
	}
	printMessage("OK. Check passed.")
	return nil
}
