package commands

import (
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory = map[string]cli.CommandFactory{
	"bot": func() (cli.Command, error) {
		return &botCommand{}, nil
	},
	"server": func() (cli.Command, error) {
		return &serverCommand{}, nil
	},
	"validate": func() (cli.Command, error) {
		return &validateCommand{}, nil
	},
}
