package slack

import (
	"bytes"
	"github.com/bivas/rivi/bot"
	"text/template"
	"time"
)

type message struct {
	Time    time.Time
	Number  int
	Title   string
	State   string
	Owner   string
	Repo    string
	Origin  string
	Targets []string
}

func buildMessage(meta bot.EventData, targets []string) *message {
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
