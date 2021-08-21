package cli

import (
	"os"

	"github.com/jessevdk/go-flags"
)

const manygramVersion = "0.2.0.dev0"

var parserFlags flags.Options = flags.HelpFlag | flags.PassDoubleDash
var parser = flags.NewParser(nil, parserFlags)

func commandHandler(command flags.Commander, args []string) error {
	if command == nil {
		parser.WriteHelp(os.Stdout)
		return nil
	}
	return command.Execute(args)
}

type profileOption struct {
	Profile struct {
		Name string `description:"Profile name" positional-arg-name:"PROFILE"`
	} `positional-args:"true" required:"false"`
}

// Run command line interface
func Run(args []string) *Error {
	parser.SubcommandsOptional = true
	parser.CommandHandler = commandHandler
	_, err := parser.ParseArgs(args)
	if err == nil {
		return nil
	}
	switch errValue := err.(type) {
	case *flags.Error:
		if errValue.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
			return nil
		}
		return newError("", errValue)
	case *Error:
		return errValue
	default:
		return newError("Unexpected error", errValue)
	}

}
