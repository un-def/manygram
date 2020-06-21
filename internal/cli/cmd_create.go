package cli

import (
	"errors"

	"github.com/un-def/manygram/internal/profile"
)

func init() {
	parser.AddCommand("create", "Create a new profile", "Create a new profile.", new(createCmd))
}

type createCmd struct {
	profileOption
	Desktop bool `short:"d" long:"desktop" description:"Also create a desktop entry"`
}

func (c *createCmd) Execute(args []string) error {
	conf, err := readConfig()
	if err != nil {
		return err
	}
	profileName := c.Profile.Name
	_, err = profile.Create(conf.ProfileDir, profileName)
	if err != nil {
		if errors.Is(err, profile.ErrInvalidName) {
			return profileNameError(profileName)
		}
		if errors.Is(err, profile.ErrAlreadyExists) {
			return newError(
				"Profile '%s' already exists. Use `manygram remove %[1]s` first if you want to recreate the profile.",
				profileName,
			)
		}
		return newError("Failed to create profile '%s'.", profileName, err)

	}
	printMessage("Profile '%s' has been created.", profileName)
	if c.Desktop {
		if err := createDesktopEntry(conf, profileName); err != nil {
			return err
		}
		printMessage("Desktop entry for profile has been created.")
	}
	return nil
}
