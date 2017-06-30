package bot

import (
	"github.com/bivas/rivi/util"
	"net/http"
)

type Bot interface {
	HandleEvent(r *http.Request)
}

type bot struct {
	configuration Configuration
}

func (b *bot) HandleEvent(r *http.Request) {
	data, process := buildFromRequest(b.configuration.GetClientConfig(), r)
	if !process {
		return
	}
	applied := make([]Rule, 0)
	for _, rule := range b.configuration.GetRules() {
		if rule.Accept(data) {
			util.Logger.Debug("Accepting rule %s for '%s'", rule.Name(), data.GetTitle())
			applied = append(applied, rule)
		}
	}
	for _, rule := range applied {
		util.Logger.Debug("Applying rule %s for '%s'", rule.Name(), data.GetTitle())
		for _, action := range rule.Actions() {
			action.Apply(b.configuration, data)
		}
	}
}

func New(configPath string) (Bot, error) {
	configuration, err := newConfiguration(configPath)
	if err != nil {
		return nil, err
	}
	return &bot{configuration: configuration}, nil
}
