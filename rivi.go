package main

import (
	"flag"
	"fmt"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/server"
	"github.com/bivas/rivi/util"
	"log"
	"os"
)

type configs []string

func (c *configs) String() string {
	return fmt.Sprintf("%s", *c)
}

func (c *configs) Set(value string) error {
	*c = append(*c, value)
	return nil
}

type botSetup struct {
	port   int
	uri    string
	config configs
}

func main() {
	logSetup()
	var setup botSetup
	flag.IntVar(&setup.port, "port", 8080, "Bot listening port")
	flag.StringVar(&setup.uri, "uri", "/", "Bot URI path")
	flag.Var(&setup.config, "config", "Bot configuration file(s)")
	flag.Parse()
	if len(setup.config) == 0 {
		util.Logger.Error("missing configuration file")
		flag.Usage()
		os.Exit(1)
	}
	run, err := bot.New(setup.config...)
	if err != nil {
		log.Fatalln("Unable to start bot handler", err)
	}
	s := server.BotServer{Port: setup.port, Uri: setup.uri, Bot: run}
	if err := s.Run(); err != nil {
		log.Fatalln("Bot exited with error. %s", err)
	}
}
