package main

import (
	"io/ioutil"
	"log"
	"os"
)

func isDebug() bool {
	return len(os.Getenv("BOT_DEBUG")) > 0
}

func logSetup() {
	log.SetOutput(os.Stderr)
	if !isDebug() {
		temp, e := ioutil.TempFile("", "rivi-bot-log.")
		if e != nil {
			panic(e)
		}
		log.Println("Log file at", temp.Name())
		log.SetOutput(temp)
	}
}
