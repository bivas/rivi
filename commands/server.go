package commands

import (
	"flag"
	"fmt"

	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/server"
	"github.com/bivas/rivi/util/log"

	"github.com/mitchellh/cli"
)

type configs []string

func (c *configs) String() string {
	return fmt.Sprintf("%s", *c)
}

func (c *configs) Set(value string) error {
	*c = append(*c, value)
	return nil
}

type serverCommand struct {
	port   int
	uri    string
	config configs
}

func (s *serverCommand) Help() string {
	return `
Usage: rivi	server [options]

	Starts rivi in server mode to listen to incoming webhooks

Options:
	-port=8080				Listen on port (default: 8080)
	-uri=/					URI path to bind POST web requests
	-config					Configuration file(s)
`
}

func (s *serverCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("server", flag.ContinueOnError)
	flagSet.IntVar(&s.port, "port", 8080, "Bot listening port")
	flagSet.StringVar(&s.uri, "uri", "/", "Bot URI path")
	flagSet.Var(&s.config, "config", "Bot configuration file(s)")
	if err := flagSet.Parse(args); err != nil {
		return 1
	}
	if len(s.config) == 0 {
		log.Error("missing configuration file")
		return cli.RunResultHelp
	}
	run, err := bot.New(s.config...)
	if err != nil {
		log.ErrorWith(log.MetaFields{log.E(err)}, "Unable to start bot handler")
		return 1
	}
	srv := server.BotServer{Port: s.port, Uri: s.uri, Bot: run}
	if err := srv.Run(); err != nil {
		log.ErrorWith(log.MetaFields{log.E(err)}, "Bot exited with error")
		return -1
	}
	return 0
}

func (s *serverCommand) Synopsis() string {
	return "start rivi server"
}
