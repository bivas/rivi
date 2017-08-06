package engine

import (
	"github.com/bivas/rivi/engine/actions"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/multistep"
)

var lrs = log.Get("engine.rule.step")

type ActionStep struct {
	action actions.Action
}

func (s *ActionStep) Run(state multistep.StateBag) multistep.StepAction {
	s.action.Apply(state)
	return multistep.ActionContinue
}

func (s *ActionStep) Cleanup(multistep.StateBag) {
}

func Run(rule Rule, state multistep.StateBag) {
	steps := make([]multistep.Step, len(rule.Actions()), len(rule.Actions()))
	for i, action := range rule.Actions() {
		steps[i] = &ActionStep{action}
	}
	meta := state.Get("data").(types.Data)
	lrs.DebugWith(
		log.MetaFields{log.F("issue", meta.GetShortName()), log.F("steps", len(steps))},
		"Applying rule %s", rule.Name())
	runner := &multistep.BasicRunner{Steps: steps}
	runner.Run(state)
}

func RunGroup(group RulesGroup, state multistep.StateBag) []string {
	meta := state.Get("data").(types.Data)
	lrs.DebugWith(
		log.MetaFields{
			log.F("issue", meta.GetShortName()),
			log.F("key", group.Key),
			log.F("rules_count", len(group.Rules)),
		}, "Processing rule group")
	applied := make([]Rule, 0)
	result := make([]string, 0)
	for _, rule := range group.Rules {
		if rule.Accept(meta) {
			lrs.DebugWith(
				log.MetaFields{
					log.F("issue", meta.GetShortName()),
					log.F("name", rule.Name()),
				}, "Accepting rule")
			applied = append(applied, rule)
			result = append(result, rule.Name())
		}
	}
	for _, rule := range applied {
		Run(rule, state)
	}
	return result
}
