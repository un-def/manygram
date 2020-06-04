package cli

func init() {
	parser.AddCommand("version", "Print manygram version", "", new(versionCmd))
}

type versionCmd struct{}

func (c *versionCmd) Execute(args []string) error {
	printMessage("manygram %s", manygramVersion)
	return nil
}
