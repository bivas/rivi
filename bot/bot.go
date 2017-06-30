package bot

import (
	"github.com/bivas/rivi/util"
	"net/http"
)

type HandledEventResult struct {
	AppliedRules []string `json:"applied_rules,omitempty"`
	Message      string   `json:"message,omitempty"`
}

type Bot interface {
	HandleEvent(r *http.Request) *HandledEventResult
}

type bot struct {
	configuration Configuration
}

func (b *bot) HandleEvent(r *http.Request) *HandledEventResult {
	result := &HandledEventResult{
		AppliedRules: []string{},
	}
	data, process := buildFromRequest(b.configuration.GetClientConfig(), r)
	if !process {
		result.Message = "Skipping rules processing (could be not supported event type)"
		return result
	}
	applied := make([]Rule, 0)
	for _, rule := range b.configuration.GetRules() {
		if rule.Accept(data) {
			util.Logger.Debug("Accepting rule %s for '%s'", rule.Name(), data.GetTitle())
			applied = append(applied, rule)
			result.AppliedRules = append(result.AppliedRules, rule.Name())
		}
	}
	for _, rule := range applied {
		util.Logger.Debug("Applying rule %s for '%s'", rule.Name(), data.GetTitle())
		for _, action := range rule.Actions() {
			action.Apply(b.configuration, data)
		}
	}
	return result
}

func New(configPath string) (Bot, error) {
	configuration, err := newConfiguration(configPath)
	if err != nil {
		return nil, err
	}
	return &bot{configuration: configuration}, nil
}
