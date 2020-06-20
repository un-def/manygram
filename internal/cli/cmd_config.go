package cli

import "github.com/jessevdk/go-flags"

var configCommand *flags.Command

func init() {
	configCommand, _ = parser.AddCommand(
		"config", "Config subcommands", "Config subcommands.",
		new(configCmd),
	)
}

type configCmd struct{}
