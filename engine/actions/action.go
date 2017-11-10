package actions

import (
	"fmt"
	"strings"

	"github.com/bivas/rivi/util/log"

	"github.com/mitchellh/multistep"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
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
				action := factory.BuildAction(config.GetStringMap(setting))
				if action != nil {
					result = append(result, action)
				}
			}
		}
	}
	return result
}

func NewCounter(name string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "rivi",
		Subsystem: "actions",
		Name:      name,
		Help:      fmt.Sprintf("Action counter for %s", name),
	})
}
