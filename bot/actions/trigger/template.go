package trigger

import (
	"bytes"
	"github.com/bivas/rivi/util"
	"io"
	"text/template"
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

var defaultTemplateBody = util.StripNonSpaceWhitespaces(`{
	"time":"{{.Time}}",
	"message":"Triggered by Rivi Bot",
	"item":{
		"repository":"{{.Owner}}/{{.Repo}}",
		"state":"{{.State}}",
		"id":{{.Number}},
		"title":"{{.Title}}"
	}
}`)

var defaultTemplate *template.Template

func init() {
	parsed, e := template.New("message").Parse(defaultTemplateBody)
	if e != nil {
		util.Logger.Error("Unable to process default template. %s", e)
	} else {
		defaultTemplate = parsed
	}
}

func processMessage(body *string, message *message) io.Reader {
	use := defaultTemplate
	if *body != "" {
		parsed, e := template.New("message").Parse(defaultTemplateBody)
		if e != nil {
			util.Logger.Error("Unable to process provided template. %s", e)
		} else {
			use = parsed
		}
	}
	var buffer bytes.Buffer
	if e := use.Execute(&buffer, message); e != nil {
		util.Logger.Error("Unable to write message to template. %s", e)
	}
	return &buffer
}
