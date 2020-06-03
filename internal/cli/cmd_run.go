package cli

import (
	"errors"
	"os"
	"path"

	"github.com/un-def/manygram/internal/config"
	"github.com/un-def/manygram/internal/profile"
	"github.com/un-def/manygram/internal/tg"
	"github.com/un-def/manygram/internal/xdg"
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
	configHome := xdg.GetConfigHome()
	configPath := path.Join(configHome, "manygram", "config.toml")
	conf, err := config.Read(configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		conf = config.Default(configPath)
		if err = conf.Write(); err != nil {
			return err
		}
	}
	profileDir := conf.ProfileDir
	if profileDir == "" {
		return errors.New("profile-dir config option is not set")
	}
	profileName := c.Profile.Name
	var prof *profile.Profile
	if c.New {
		if prof, err = profile.New(profileDir, profileName); err != nil {
			return err
		}
	} else if prof, err = profile.Read(profileDir, profileName); err != nil {
		return err
	}
	return tg.Run(conf.BinPath, prof.Path, args, c.Wait)
}
