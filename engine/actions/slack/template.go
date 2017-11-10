package slack

import (
	"bytes"
	"text/template"
	"time"

	"github.com/bivas/rivi/types"
)

type message struct {
	Time    time.Time
	Number  int
	Title   string
	State   string
	Owner   string
	Repo    string
	Origin  types.Origin
	Targets []string
}

func buildMessage(meta types.Data, targets []string) *message {
	return &message{
		Time:    time.Now(),
		Number:  meta.GetNumber(),
		Title:   meta.GetTitle(),
		State:   meta.GetState(),
		Owner:   meta.GetOwner(),
		Repo:    meta.GetRepo(),
		Origin:  meta.GetOrigin(),
		Targets: targets,
	}
}

func serializeMessage(t *template.Template, message *message) (string, error) {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, message); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
