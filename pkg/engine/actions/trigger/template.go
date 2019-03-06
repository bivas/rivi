package trigger

import (
	"bytes"
	"io"
	"text/template"
	"time"

	"github.com/bivas/rivi/pkg/types"
	"github.com/bivas/rivi/pkg/util"
	"github.com/bivas/rivi/pkg/util/log"
)

type message struct {
	Time   time.Time
	Number int
	Title  string
	State  string
	Owner  string
	Repo   string
	Origin types.Origin
}

var defaultTemplateBody = util.StripNonSpaceWhitespaces(`{
	"time":"{{.Time}}",
	"message":"Triggered by Rivi",
	"item":{
		"repository":"{{.Owner}}/{{.Repo}}",
		"state":"{{.State}}",
		"id":{{.Number}},
		"title":"{{.Title}}"
	}
}`)

var (
	defaultTemplate *template.Template
	logger          = log.Get("trigger.template")
)

func init() {
	parsed, e := template.New("message").Parse(defaultTemplateBody)
	if e != nil {
		logger.ErrorWith(log.MetaFields{log.E(e)}, "Unable to process default template")
	} else {
		defaultTemplate = parsed
	}
}

func processMessage(body *string, message *message) io.Reader {
	use := defaultTemplate
	if *body != "" {
		parsed, e := template.New("message").Parse(*body)
		if e != nil {
			logger.ErrorWith(log.MetaFields{log.E(e)}, "Unable to process provided template")
		} else {
			use = parsed
		}
	}
	var buffer bytes.Buffer
	if e := use.Execute(&buffer, message); e != nil {
		logger.ErrorWith(log.MetaFields{log.E(e)}, "Unable to write message to template")
	}
	return &buffer
}
