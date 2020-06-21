package cli

import (
	"errors"

	"github.com/un-def/manygram/internal/profile"
)

func init() {
	parser.AddCommand("remove", "Remove the profile", "Remove the profile.", new(removeCmd))
}

type removeCmd struct {
	profileOption
	Desktop bool `short:"d" long:"desktop" description:"Also remove the desktop entry"`
}

func (c *removeCmd) Execute(args []string) error {
	conf, err := readConfig()
	if err != nil {
		return err
	}
	profileName := c.Profile.Name
	if err = profile.Remove(conf.ProfileDir, profileName); err != nil {
		if errors.Is(err, profile.ErrInvalidName) {
			return profileNameError(profileName)
		}
		if errors.Is(err, profile.ErrNotExist) {
			return newError("Profile '%s' does not exist.", profileName)
		}
		return newError("Failed to remove profile '%s'.", profileName, err)
	}
	printMessage("Profile '%s' has been removed.", profileName)
	if c.Desktop {
		if err := removeDesktopEntry(profileName); err != nil {
			return err
		}
		printMessage("Desktop entry for profile has been removed.")
	}
	return nil
}
