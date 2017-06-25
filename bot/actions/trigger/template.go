package trigger

import (
	"github.com/bivas/rivi/util"
	"time"
)

type message struct {
	Time   time.Time
	Number int
	Title  string
	State  string
	Owner  string
	Repo   string
	Origin string
}

var defaultTemplate = util.StripNonSpaceWhitespaces(`{
	"time":"{{.Time}}",
	"message":"Triggered by Rivi Bot",
	"item":{
		"repository":"{{.Owner}}/{{.Repo}}",
		"state":"{{.State}}",
		"id":{{.Number}},
		"title":"{{.Title}}"
	}
}`)
