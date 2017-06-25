package main

import (
	"flag"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/server"
	"github.com/bivas/rivi/util"
	"log"
	"os"
)

func main() {
	logSetup()
	port := flag.Int("port", 8080, "Bot listening port")
	uri := flag.String("uri", "/", "Bot URI path")
	config := flag.String("config", "", "Bot configuration file")
	flag.Parse()
	if *config == "" {
		util.Logger.Error("missing configuration file")
		flag.Usage()
		os.Exit(1)
	}
	run, err := bot.New(*config)
	if err != nil {
		log.Fatalln("Unable to start bot handler", err)
	}
	s := server.BotServer{Port: *port, Uri: *uri, Bot: run}
	s.Run()
}
