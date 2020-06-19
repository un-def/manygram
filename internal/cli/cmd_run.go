package cli

import (
	"errors"
	"os"

	"github.com/un-def/manygram/internal/config"
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
	New  bool `short:"n" long:"new" description:"Create a new profile"`
	Wait bool `short:"w" long:"wait" description:"Wait for child process to terminate"`
}

func (c *runCmd) Execute(args []string) error {
	configPath := getConfigPath()
	conf, err := config.Read(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return newError(
				"Config %s not found. Run `manygram config` to create a new config.",
				configPath, err,
			)
		}
		return newError("Failed to read config %s", configPath, err)
	}

	telegram, err := tg.Executable(conf.ExecPath)
	if err != nil {
		return newError("Failed to locale Telegram Desktop executable. Check `exec-path` config parameter.", err)
	}

	profileDir := conf.ProfileDir
	if profileDir == "" {
		return newError("`profile-dir` config parameter is not set.")
	}

	profileName := c.Profile.Name
	var prof *profile.Profile
	if c.New {
		prof, err = profile.New(profileDir, profileName)
	} else {
		prof, err = profile.Read(profileDir, profileName)
	}
	if err != nil {
		if errors.Is(err, profile.ErrInvalidName) {
			return profileNameError(profileName)
		}
		if errors.Is(err, profile.ErrNotExist) {
			return newError("Profile '%s' does not exist. Use `--new` flag to create a new one.", profileName)
		}
		if errors.Is(err, profile.ErrAlreadyExists) {
			return newError(
				"Profile '%s' already exists. Use `manygram remove %[1]s` first if you want to recreate the profile.",
				profileName,
			)
		}
		return newError("Failed to read or create profile '%s'.", profileName, err)
	}
	return telegram.Run(prof.Path, args, c.Wait)
}
