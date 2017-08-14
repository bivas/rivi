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
Usage: rivi	server [options] CONFIGURATION_FILE(S)...

	Starts rivi in server mode and listen to incoming webhooks

Options:
	-port=8080				Listen on port (default: 8080)
	-uri=/					URI path (default: "/")
`
}

func (s *serverCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("server", flag.ContinueOnError)
	flagSet.IntVar(&s.port, "port", 8080, "Runner listening port")
	flagSet.StringVar(&s.uri, "uri", "/", "Runner URI path")
	if err := flagSet.Parse(args); err != nil {
		return cli.RunResultHelp
	}
	if len(flagSet.Args()) == 0 {
		log.Error("missing configuration file")
		return cli.RunResultHelp
	}
	run, err := runner.New(flagSet.Args()...)
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
