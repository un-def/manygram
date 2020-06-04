package main

import (
	"fmt"
	"os"

	"github.com/un-def/manygram/internal/cli"
)

func main() {
	err := cli.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.NiceError())
		os.Exit(1)
	}
}
