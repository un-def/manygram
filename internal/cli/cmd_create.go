package cli

import (
	"errors"

	"github.com/un-def/manygram/internal/profile"
)

func init() {
	parser.AddCommand("create", "Create a new profile", "Create a new profile", new(createCmd))
}

type createCmd struct {
	Profile struct {
		Name string `description:"Profile name" positional-arg-name:"PROFILE"`
	} `positional-args:"true" required:"true"`
}

func (c *createCmd) Execute(args []string) error {
	profileDir, err := getProfileDirParameter(nil)
	if err != nil {
		return err
	}
	profileName := c.Profile.Name
	_, err = profile.New(profileDir, profileName)
	if err == nil {
		printMessage("Profile '%s' has been created.", profileName)
		return nil
	}
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
