package engine

import (
	"fmt"
	"sort"

	"github.com/bivas/rivi/engine/actions"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/multistep"
	"github.com/spf13/viper"
)

var lr = log.Get("engine.rule")

type Rule interface {
	Name() string
	Order() int
	Accept(meta types.Data) bool
	Actions() []actions.Action
}

type RulesByConditionOrder []Rule

func (r RulesByConditionOrder) Len() int {
	return len(r)
}

func (r RulesByConditionOrder) Less(i, j int) bool {
	return r[i].Order() < r[j].Order()
}

func (r RulesByConditionOrder) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type rule struct {
	name      string
	condition Condition
	actions   []actions.Action
}

func (r *rule) Name() string {
	return r.name
}

func (r *rule) Order() int {
	return r.condition.Order
}

func (r *rule) String() string {
	format := `{
		name: %s
		condition: %+v
		actions: %s
	}`
	as := "{\n"
	for index, action := range r.actions {
		as += fmt.Sprintf("			%d { %s }\n", index+1, action)
	}
	as += "		}"
	return fmt.Sprintf(format, r.name, r.condition, as)
}

func (r *rule) Accept(meta types.Data) bool {
	accept := r.condition.Match(meta)
	if !accept {
		lr.DebugWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "Skipping rule '%s'", r.name)
	}
	return accept
}

func (r *rule) Actions() []actions.Action {
	return r.actions
}

type RulesGroup struct {
	Key   int
	Rules []Rule
}

type rulesGroups []RulesGroup

func (r rulesGroups) Len() int {
	return len(r)
}

func (r rulesGroups) Less(i, j int) bool {
	return r[i].Key < r[j].Key
}

func (r rulesGroups) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func GroupByRuleOrder(rules []Rule) []RulesGroup {
	groupIndexes := make(map[int]RulesGroup)
	for _, rule := range rules {
		key := rule.Order()
		group, exists := groupIndexes[key]
		if !exists {
			group = RulesGroup{key, make([]Rule, 0)}
		}
		group.Rules = append(group.Rules, rule)
		groupIndexes[key] = group
	}
	lr.Debug("%d rules are grouped to %d rule groups", len(rules), len(groupIndexes))
	groupsResult := make([]RulesGroup, 0)
	for _, group := range groupIndexes {
		groupsResult = append(groupsResult, group)
	}
	sort.Sort(rulesGroups(groupsResult))
	return groupsResult
}

func NewRule(name string, config *viper.Viper) Rule {
	return &rule{
		name:      name,
		condition: buildConditionFromConfiguration(config),
		actions:   actions.BuildActionsFromConfiguration(config),
	}
}

func ProcessRules(rules []Rule, state multistep.StateBag) []string {
	groups := GroupByRuleOrder(rules)
	result := make([]string, 0)
	for _, group := range groups {
		result = append(result, RunGroup(group, state)...)
	}
	return result
}
