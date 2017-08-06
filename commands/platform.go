package commands

import (
	"flag"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/server"
	"github.com/bivas/rivi/util/log"

	"github.com/mitchellh/cli"
)

type platformCommand struct {
	port int
	uri  string
}

func (s *platformCommand) Help() string {
	return `
Usage: rivi	platform [options]

	Starts rivi in platform mode and listen to incoming webhooks

Options:
	-port=8080				Listen on port (default: 8080)
	-uri=/					URI path (default: "/")
`
}

func (s *platformCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("server", flag.ContinueOnError)
	flagSet.IntVar(&s.port, "port", 8080, "Bot listening port")
	flagSet.StringVar(&s.uri, "uri", "/", "Bot URI path")
	if err := flagSet.Parse(args); err != nil {
		return cli.RunResultHelp
	}
	run, err := bot.NewHookListener()
	if err != nil {
		log.ErrorWith(log.MetaFields{log.E(err)}, "Unable to start bot handler")
		return 1
	}
	log.Info("Rivi is ready")
	srv := server.BotServer{Port: s.port, Uri: s.uri, Bot: run}
	if err := srv.Run(); err != nil {
		log.ErrorWith(log.MetaFields{log.E(err)}, "Bot exited with error")
		return -1
	}
	return 0
}

func (s *platformCommand) Synopsis() string {
	return "start rivi server"
}
