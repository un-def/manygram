package cli

import (
	"os"

	"github.com/jessevdk/go-flags"
)

var parserFlags flags.Options = flags.HelpFlag | flags.PassDoubleDash
var parser = flags.NewParser(nil, parserFlags)

func commandHandler(command flags.Commander, args []string) error {
	if command == nil {
		parser.WriteHelp(os.Stdout)
		return nil
	}
	return command.Execute(args)
}

// Run command line interface
func Run(args []string) *Error {
	parser.SubcommandsOptional = true
	parser.CommandHandler = commandHandler
	if _, err := parser.ParseArgs(args); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
			return nil
		}
		if cliErr, ok := err.(*Error); ok {
			return cliErr
		}
		return newError("Unexpected error", err)
	}
	return nil
}
