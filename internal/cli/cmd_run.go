package cli

import (
	"errors"

	"github.com/un-def/manygram/internal/profile"
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
	Profile struct {
		Name string `description:"Profile name" positional-arg-name:"PROFILE"`
	} `positional-args:"true" required:"true"`
	Wait bool `short:"w" long:"wait" description:"Wait for child process to terminate"`
}

func (c *runCmd) Execute(args []string) error {
	conf, err := readConfig()
	if err != nil {
		return err
	}
	telegram, err := tg.Executable(conf.ExecPath)
	if err != nil {
		return newError("Failed to locale Telegram Desktop executable. Check `exec-path` config parameter.", err)
	}
	profileName := c.Profile.Name
	prof, err := profile.Read(conf.ProfileDir, profileName)
	if err != nil {
		if errors.Is(err, profile.ErrInvalidName) {
			return profileNameError(profileName)
		}
		if errors.Is(err, profile.ErrNotExist) {
			return newError("Profile '%s' does not exist. Use `manygram create %[1]s` to create a new one.", profileName)
		}
		return newError("Failed to read profile '%s'.", profileName, err)
	}
	return telegram.Run(prof.Path, args, c.Wait)
}
