package actions

import (
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/multistep"
	"github.com/spf13/viper"
	"strings"
)

var la = log.Get("actions")

type Action interface {
	Apply(multistep.StateBag)
}

type ActionFactory interface {
	BuildAction(config map[string]interface{}) Action
}

var registry map[string]ActionFactory = make(map[string]ActionFactory)
var supportedActions []string = make([]string, 0)

func RegisterAction(kind string, action ActionFactory) {
	search := strings.ToLower(kind)
	_, exists := registry[search]
	if exists {
		la.Error("action %s exists!", kind)
	} else {
		la.Debug("registering action %s", kind)
		registry[search] = action
		supportedActions = append(supportedActions, kind)
	}
}

func BuildActionsFromConfiguration(config *viper.Viper) []Action {
	result := make([]Action, 0)
	for setting := range config.AllSettings() {
		if setting == "condition" {
			continue
		}
		for _, support := range supportedActions {
			if setting == support {
				factory := registry[setting]
				result = append(result, factory.BuildAction(config.GetStringMap(setting)))
			}
		}
	}
	return result
}
