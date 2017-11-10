package main

import (
	"os"

	"github.com/bivas/rivi/commands"

	"github.com/mitchellh/cli"
)

func main() {
	logSetup()
	c := cli.CLI{
		Name:         "rivi",
		Autocomplete: true,
		Args:         os.Args[1:],
		Commands:     commands.Commands,
		HelpFunc:     cli.BasicHelpFunc("rivi"),
		HelpWriter:   os.Stdout,
	}
	exitCode, _ := c.Run()
	os.Exit(exitCode)
}
