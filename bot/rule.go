package bot

import (
	"fmt"
	"github.com/bivas/rivi/util"
	"github.com/spf13/viper"
	"strings"
)

type Action interface {
	Apply(config Configuration, meta EventData)
}

type Rule interface {
	Name() string
	Order() int
	Accept(meta EventData) bool
	Actions() []Action
}

type rulesByConditionOrder []Rule

func (r rulesByConditionOrder) Len() int {
	return len(r)
}

func (r rulesByConditionOrder) Less(i, j int) bool {
	return r[i].Order() < r[j].Order()
}

func (r rulesByConditionOrder) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type rule struct {
	name      string
	condition Condition
	actions   []Action
}

func (r *rule) Name() string {
	return r.name
}

func (r *rule) Order() int {
	return r.condition.Order
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

func (r *rule) Actions() []Action {
	return r.actions
}

type ActionFactory interface {
	BuildAction(config map[string]interface{}) Action
}

var actions map[string]ActionFactory = make(map[string]ActionFactory)
var supportedActions []string = make([]string, 0)

func RegisterAction(kind string, action ActionFactory) {
	search := strings.ToLower(kind)
	_, exists := actions[search]
	if exists {
		util.Logger.Error("action %s exists!", kind)
	} else {
		util.Logger.Debug("registering action %s", kind)
		actions[search] = action
		supportedActions = append(supportedActions, kind)
	}
}

func buildActionsFromConfiguration(config *viper.Viper) []Action {
	result := make([]Action, 0)
	for setting := range config.AllSettings() {
		if setting == "condition" {
			continue
		}
		for _, support := range supportedActions {
			if setting == support {
				factory := actions[setting]
				result = append(result, factory.BuildAction(config.GetStringMap(setting)))
			}
		}
	}
	return result
}
