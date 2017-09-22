package commands

import (
	"flag"

	"github.com/bivas/rivi/runner"
	"github.com/bivas/rivi/server"
	"github.com/bivas/rivi/util/log"

	"github.com/mitchellh/cli"
)

type serverCommand struct {
	port int
	uri  string
}

func (s *serverCommand) Help() string {
	return `
Usage: rivi	server [options] [config]

	Starts rivi in server mode and listen to incoming webhooks

Options:
	-port=8080				Listen on port
	-uri=/					URI path
`
}

func (s *serverCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("server", flag.ContinueOnError)
	flagSet.IntVar(&s.port, "port", 8080, "Runner listening port")
	flagSet.StringVar(&s.uri, "uri", "/", "Runner URI path")
	if err := flagSet.Parse(args); err != nil {
		return cli.RunResultHelp
	}

	configFile := ""
	switch flagSet.NArg() {
	case 0:
		// nothing
	case 1:
		configFile = flagSet.Args()[0]
	default:
		log.Error("Too many args provided (expected only 1)")
		return cli.RunResultHelp
	}

	run, err := runner.NewHookListener(configFile)
	if err != nil {
		log.ErrorWith(log.MetaFields{log.E(err)}, "Unable to start runner handler")
		return 1
	}
	log.Info("Rivi is ready")
	srv := server.BotServer{Port: s.port, Uri: s.uri, Runner: run}
	if err := srv.Run(); err != nil {
		log.ErrorWith(log.MetaFields{log.E(err)}, "Runner exited with error")
		return -1
	}
	return 0
}

func (s *serverCommand) Synopsis() string {
	return "start rivi server"
}
