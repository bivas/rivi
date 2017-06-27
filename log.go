package main

import (
	"io/ioutil"
	"log"
	"os"
)

func isDebug() bool {
	if len(os.Getenv("BOT_DEBUG")) > 0 {
		return true
	}
	return false
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
