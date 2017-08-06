package commands

import (
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory = map[string]cli.CommandFactory{
	"server": func() (cli.Command, error) {
		return &serverCommand{}, nil
	},
	"platform": func() (cli.Command, error) {
		return &platformCommand{}, nil
	},
	"validate": func() (cli.Command, error) {
		return &validateCommand{}, nil
	},
}
