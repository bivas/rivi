package bot

import (
	"fmt"
	"github.com/bivas/rivi/util"
	"github.com/spf13/viper"
)

type Action interface {
	Apply(config Configuration, meta EventData)
}

type Rule interface {
	Name() string
	Accept(meta EventData) bool
	Action() Action
}

type rule struct {
	name      string
	condition Condition
	action    Action
}

func (r *rule) Name() string {
	return r.name
}

func (r *rule) String() string {
	return fmt.Sprintf("%#v", r)
}

func (r *rule) Accept(meta EventData) bool {
	accept := r.condition.Match(meta)
	if !accept {
		util.Logger.Debug("Skipping rule '%s'", r.name)
	}
	return accept
}

func (r *rule) Action() Action {
	return r.action
}

type ActionFactory interface {
	BuildAction(config map[string]interface{}) Action
}

var actions map[string]ActionFactory = make(map[string]ActionFactory)
var supportedActions []string = make([]string, 0)

func RegisterAction(kind string, action ActionFactory) {
	actions[kind] = action
	supportedActions = append(supportedActions, kind)
	util.Logger.Debug("running with support for %s", kind)
}

func buildActionFromConfiguration(config *viper.Viper) Action {
	for setting := range config.AllSettings() {
		if setting == "condition" {
			continue
		}
		for _, support := range supportedActions {
			if setting == support {
				factory := actions[setting]
				return factory.BuildAction(config.GetStringMap(setting))
			}
		}
	}
	return nil
}
