package engine

import (
	"testing"

	"sort"

	"github.com/bivas/rivi/engine/actions"
	"github.com/bivas/rivi/types"
	"github.com/mitchellh/multistep"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testRule struct {
	NameValue  string
	OrderValue int
}

func (r *testRule) Name() string {
	return r.NameValue
}

func (r *testRule) Order() int {
	return r.OrderValue
}

func (r *testRule) Accept(meta types.Data) bool {
	panic("implement me")
}

func (r *testRule) Actions() []actions.Action {
	panic("implement me")
}

func TestRulesByConditionOrder(t *testing.T) {
	rules := []Rule{
		&testRule{OrderValue: 5},
		&testRule{OrderValue: 4},
		&testRule{OrderValue: 3},
		&testRule{OrderValue: 2},
	}
	sort.Sort(RulesByConditionOrder(rules))
	assert.Equal(t, rules[0].Order(), 2, "1st")
	assert.Equal(t, rules[1].Order(), 3, "2nd")
	assert.Equal(t, rules[2].Order(), 4, "3rd")
	assert.Equal(t, rules[3].Order(), 5, "4th")
}

func TestGroupByRuleOrder(t *testing.T) {
	rules := []Rule{
		&testRule{OrderValue: 5, NameValue: "rule1"},
		&testRule{OrderValue: 5, NameValue: "rule2"},
		&testRule{OrderValue: 3, NameValue: "rule3"},
		&testRule{OrderValue: 3, NameValue: "rule4"},
	}
	groups := GroupByRuleOrder(rules)
	assert.Len(t, groups, 2, "more than 2")
}

func TestProcessRules(t *testing.T) {
	emptyViper := viper.New()
	names := []string{"rule1", "rule2", "rule3"}
	rules := []Rule{}
	for _, name := range names {
		rules = append(rules, NewRule(name, emptyViper))
	}
	state := &multistep.BasicStateBag{}
	state.Put("data", &mockData{})
	result := ProcessRules(rules, state)
	assert.Len(t, result, 3, "accept all 3")
}
