package cli

import "github.com/jessevdk/go-flags"

var desktopCommand *flags.Command

func init() {
	desktopCommand, _ = parser.AddCommand(
		"desktop", "Desktop entries subcommands", "Desktop entries subcommands.",
		new(desktopCmd),
	)
}

type desktopCmd struct{}
