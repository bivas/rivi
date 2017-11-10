package commands

import (
	"flag"
	"github.com/bivas/rivi/runner"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/cli"
)

type validateCommand struct {
}

func (v *validateCommand) Help() string {
	return `
Usage: rivi	validate CONFIGURATION_FILE(S)...

	Validate inputted configuration file(s)
`
}

func (v *validateCommand) Run(args []string) int {
	f := flag.NewFlagSet("validate", flag.ContinueOnError)
	if err := f.Parse(args); err != nil {
		return cli.RunResultHelp
	}
	if len(f.Args()) == 0 {
		log.Error("Missing configuration files to validate")
		return cli.RunResultHelp
	}
	hadError := false
	for _, file := range f.Args() {
		_, err := runner.New(file)
		if err != nil {
			log.ErrorWith(log.MetaFields{log.F("config", file), log.E(err)}, "Config file failed")
			hadError = true
		} else {
			log.InfoWith(log.MetaFields{log.F("config", file)}, "Config file passed")
		}
	}
	if hadError {
		return -1
	}
	return 0
}

func (v *validateCommand) Synopsis() string {
	return "validate configuration file"
}
