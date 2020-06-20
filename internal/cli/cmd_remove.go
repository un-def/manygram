package cli

import (
	"errors"

	"github.com/un-def/manygram/internal/profile"
)

func init() {
	parser.AddCommand("remove", "Remove the profile", "Remove the profile", new(removeCmd))
}

type removeCmd struct {
	Profile struct {
		Name string `description:"Profile name" positional-arg-name:"PROFILE"`
	} `positional-args:"true" required:"true"`
}

func (c *removeCmd) Execute(args []string) error {
	conf, err := readConfig()
	if err != nil {
		return err
	}
	profileName := c.Profile.Name
	if err = profile.Delete(conf.ProfileDir, profileName); err == nil {
		printMessage("Profile '%s' has been removed.", profileName)
		return nil
	}
	if errors.Is(err, profile.ErrInvalidName) {
		return profileNameError(profileName)
	}
	if errors.Is(err, profile.ErrNotExist) {
		return newError("Profile '%s' does not exist.", profileName)
	}
	return newError("Failed to remove profile '%s'.", profileName, err)
}
