package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/bivas/rivi/cmd/rivi/app"

	"github.com/mitchellh/cli"

	// actions
	_ "github.com/bivas/rivi/pkg/engine/actions/autoassign"
	_ "github.com/bivas/rivi/pkg/engine/actions/automerge"
	_ "github.com/bivas/rivi/pkg/engine/actions/commenter"
	_ "github.com/bivas/rivi/pkg/engine/actions/labeler"
	_ "github.com/bivas/rivi/pkg/engine/actions/locker"
	_ "github.com/bivas/rivi/pkg/engine/actions/sizing"
	_ "github.com/bivas/rivi/pkg/engine/actions/slack"
	_ "github.com/bivas/rivi/pkg/engine/actions/status"
	_ "github.com/bivas/rivi/pkg/engine/actions/trigger"

	// connectors
	_ "github.com/bivas/rivi/pkg/connectors/github"
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

func main() {
	logSetup()
	c := cli.CLI{
		Name:         "rivi",
		Autocomplete: true,
		Args:         os.Args[1:],
		Commands:     app.Commands,
		HelpFunc:     cli.BasicHelpFunc("rivi"),
		HelpWriter:   os.Stdout,
	}
	exitCode, _ := c.Run()
	os.Exit(exitCode)
}
