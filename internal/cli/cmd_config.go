package cli

import (
	"errors"
	"os"
	"strings"

	"github.com/un-def/manygram/internal/config"
	"github.com/un-def/manygram/internal/tg"
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
		conf = config.Default(configPath)

		if err = conf.Write(); err != nil {
			return newError("Failed to write config.", err)
		}
		printMessage("Done.")
		return nil
	}

	printMessage("Config %s found. Checking.", configPath)
	telegram, err := tg.Executable(conf.ExecPath)
	if err != nil {
		return newError(`Check error at "exec-path".`, err)
	}

	execMsg := []string{telegram.Path}
	if telegram.FullPath != telegram.Path {
		execMsg = append(execMsg, telegram.FullPath)
	}
	if telegram.RealPath != telegram.FullPath {
		execMsg = append(execMsg, telegram.RealPath)
	}

	printMessage("Telegram Desktop executable: %s", strings.Join(execMsg, " -> "))
	printMessage("Profile directory: %s", conf.ProfileDir)
	printMessage("OK. Check passed.")
	return nil
}
