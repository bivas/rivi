package main

import (
	"io/ioutil"
	"log"
	"os"
)

func logSetup() {
	log.SetOutput(os.Stderr)
	if len(os.Getenv("BOT_DEBUG")) == 0 {
		temp, e := ioutil.TempFile("", "review-bot-log.")
		if e != nil {
			panic(e)
		}
		log.Println("Log file at", temp.Name())
		log.SetOutput(temp)
	}
}
