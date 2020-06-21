package cli

func init() {
	desktopCommand.AddCommand(
		"create", "Create a new desktop entry", "Create a new desktop entry.",
		new(desktopCreateCmd),
	)
}

type desktopCreateCmd struct {
	profileOption
}

func (c *desktopCreateCmd) Execute(args []string) error {
	profileName := c.Profile.Name
	if err := createDesktopEntry(nil, profileName); err != nil {
		return err
	}
	printMessage("Desktop entry for profile '%s' has been created.", profileName)
	return nil
}
