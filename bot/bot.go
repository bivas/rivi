package bot

import (
	"fmt"
	"github.com/bivas/rivi/util"
	"net/http"
	"path/filepath"
	"strings"
)

type HandledEventResult struct {
	AppliedRules []string `json:"applied_rules,omitempty"`
	Message      string   `json:"message,omitempty"`
}

type Bot interface {
	HandleEvent(r *http.Request) *HandledEventResult
}

type bot struct {
	defaultNamespace string
	configurations   map[string]Configuration
}

func (b *bot) getCurrentConfiguration(namespace string) (Configuration, error) {
	if namespace == "" {
		namespace = b.defaultNamespace
	}
	configuration, exists := b.configurations[namespace]
	if !exists {
		util.Logger.Warning("Request for namespace '%s' matched nothing", namespace)
		return nil, fmt.Errorf("Request for namespace '%s' matched nothing", namespace)
	}
	return configuration, nil
}

func (b *bot) HandleEvent(r *http.Request) *HandledEventResult {
	result := &HandledEventResult{
		AppliedRules: []string{},
	}
	workingConfiguration, err := b.getCurrentConfiguration(r.URL.Query().Get("namespace"))
	if err != nil {
		result.Message = err.Error()
		return result
	}
	data, process := buildFromRequest(workingConfiguration.GetClientConfig(), r)
	if !process {
		result.Message = "Skipping rules processing (could be not supported event type)"
		return result
	}
	applied := make([]Rule, 0)
	for _, rule := range workingConfiguration.GetRules() {
		if rule.Accept(data) {
			util.Logger.Debug("Accepting rule %s for '%s'", rule.Name(), data.GetTitle())
			applied = append(applied, rule)
			result.AppliedRules = append(result.AppliedRules, rule.Name())
		}
	}
	for _, rule := range applied {
		util.Logger.Debug("Applying rule %s for '%s'", rule.Name(), data.GetTitle())
		for _, action := range rule.Actions() {
			action.Apply(workingConfiguration, data)
		}
	}
	return result
}

func New(configPaths ...string) (Bot, error) {
	b := &bot{configurations: make(map[string]Configuration)}
	for index, configPath := range configPaths {
		baseConfigPath := filepath.Base(configPath)
		namespace := strings.TrimSuffix(baseConfigPath, filepath.Ext(baseConfigPath))
		util.Logger.Debug("Loading configuration for namespace '%s'", namespace)
		if index == 0 {
			b.defaultNamespace = namespace
		}
		configuration, err := newConfiguration(configPath)
		if err != nil {
			return nil, fmt.Errorf("Reading %s caused an error. %s", configPath, err)
		}
		b.configurations[namespace] = configuration
	}
	if len(b.configurations) == 0 {
		return nil, fmt.Errorf("Bot has no readable configuration!")
	}
	util.Logger.Debug("Bot is ready %+v", *b)
	return b, nil
}
