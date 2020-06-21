package cli

func init() {
	desktopCommand.AddCommand(
		"remove", "Remove the desktop entry", "Remove the desktop entry.",
		new(desktopRemoveCmd),
	)
}

type desktopRemoveCmd struct {
	profileOption
}

func (c *desktopRemoveCmd) Execute(args []string) error {
	profileName := c.Profile.Name
	err := removeDesktopEntry(profileName)
	if err != nil {
		return err
	}
	printMessage("Desktop entry for profile '%s' has been removed.", profileName)
	return nil
}
