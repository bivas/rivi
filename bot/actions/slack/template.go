package slack

import (
	"bytes"
	"github.com/bivas/rivi/bot"
	"text/template"
	"time"
)

type message struct {
	Time      time.Time
	Number    int
	Title     string
	State     string
	Owner     string
	Repo      string
	Origin    string
	SlackUser string
}

func buildFromEventData(meta bot.EventData, slackUser string) *message {
	return &message{
		Time:      time.Now(),
		Number:    meta.GetNumber(),
		Title:     meta.GetTitle(),
		State:     meta.GetState(),
		Owner:     meta.GetOwner(),
		Repo:      meta.GetRepo(),
		Origin:    meta.GetOrigin(),
		SlackUser: slackUser,
	}
}

func buildMessage(t *template.Template, message *message) (string, error) {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, message); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
