package cli

import (
	"github.com/un-def/manygram/internal/tg"
)

func init() {
	parser.AddCommand("run", "Run Telegram Desktop", `
		Run Telegram Desktop with specified profile.
		Any additional arguments after double dash delimiter '--'
		will be passed to Telegram Desktop executable.
	`, new(runCmd))
}

type runCmd struct {
	profileOption
	Wait bool `short:"w" long:"wait" description:"Wait for child process to terminate"`
}

func (c *runCmd) Execute(args []string) error {
	conf, err := readConfig()
	if err != nil {
		return err
	}
	telegram, err := tg.Executable(conf.ExecPath, conf.ExecArgs)
	if err != nil {
		return newError("Failed to locate Telegram Desktop executable. Check `exec-path` config parameter.", err)
	}
	profileName := c.Profile.Name
	prof, err := readProfile(conf.ProfileDir, profileName)
	if err != nil {
		return err
	}
	return telegram.Run(prof.Path, args, c.Wait)
}
