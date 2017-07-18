package bot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bivas/rivi/util/log"
	"github.com/spf13/viper"
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
		log.DebugWith(log.MetaFields{{"issue", meta.GetShortName()}}, "Skipping rule '%s'", r.name)
	}
	return accept
}

func (r *rule) Actions() []Action {
	return r.actions
}

type rulesGroup struct {
	key   int
	rules []Rule
}

type rulesGroups []rulesGroup

func (r rulesGroups) Len() int {
	return len(r)
}

func (r rulesGroups) Less(i, j int) bool {
	return r[i].key < r[j].key
}

func (r rulesGroups) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func groupByRuleOrder(rules []Rule) []rulesGroup {
	groupIndexes := make(map[int]rulesGroup)
	for _, rule := range rules {
		key := rule.Order()
		group, exists := groupIndexes[key]
		if !exists {
			group = rulesGroup{key, make([]Rule, 0)}
		}
		group.rules = append(group.rules, rule)
		groupIndexes[key] = group
	}
	log.Debug("%d Rules are grouped to %d rule groups", len(rules), len(groupIndexes))
	groupsResult := make([]rulesGroup, 0)
	for _, group := range groupIndexes {
		groupsResult = append(groupsResult, group)
	}
	sort.Sort(rulesGroups(groupsResult))
	return groupsResult
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
		log.Error("action %s exists!", kind)
	} else {
		log.Debug("registering action %s", kind)
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
