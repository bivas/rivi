package bot

import (
	"fmt"
	"github.com/bivas/rivi/util"
	"github.com/spf13/viper"
	"regexp"
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
	return fmt.Sprintf("%T[name: %s,condition: %+v,action: %+v]", r, r.name, r.condition, r.action)
}

func (r *rule) checkIfLabeled(meta *EventData) bool {
	accept := false
	if len(r.condition.IfLabeled) == 0 {
		accept = true
	} else {
		for _, check := range r.condition.IfLabeled {
			for _, label := range (*meta).GetLabels() {
				accept = accept || check == label
			}
		}
	}
	return accept
}

func (r *rule) checkPattern(meta *EventData) bool {
	if r.condition.Filter.Pattern == "" {
		return true
	} else {
		for _, check := range (*meta).GetFileNames() {
			matched, e := regexp.MatchString(r.condition.Filter.Pattern, check)
			if e != nil {
				util.Logger.Debug("Error checking filter %s", e)
			} else if matched {
				return true
			}
		}
	}
	return false
}

func (r *rule) checkExt(meta *EventData) bool {
	if r.condition.Filter.Extension == "" {
		return true
	} else {
		for _, check := range (*meta).GetFileExtensions() {
			accept := r.condition.Filter.Extension == check
			if accept {
				return true
			}
		}
	}
	return false
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
