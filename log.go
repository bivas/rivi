package main

import (
	"io/ioutil"
	"log"
	"os"
)

func isDebug() bool {
	return len(os.Getenv("RIVI_DEBUG")) > 0
}

func logSetup() {
	log.SetOutput(os.Stderr)
	if !isDebug() {
		temp, e := ioutil.TempFile("", "rivi-log.")
		if e != nil {
			panic(e)
		}
		log.Println("Log file at", temp.Name())
		log.SetOutput(temp)
	}
}
