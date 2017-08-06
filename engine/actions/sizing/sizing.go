package sizing

import (
	"sort"

	"github.com/bivas/rivi/engine/actions"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/multistep"
)

type action struct {
	items          rules
	possibleLabels []string
	logger         log.Logger
}

func (a *action) updatePossibleLabels() {
	set := util.StringSet{}
	for _, item := range a.items {
		set.Add(item.Label)
	}
	a.possibleLabels = set.Values()
}

func (a *action) findMatchedLabel(meta types.Data) (*sizingRule, string, bool) {
	changedFiles := meta.GetChangedFiles()
	add, del := meta.GetChanges()
	changes := add + del
	defaultLabel := ""
	defaultExists := false
	var defaultRule sizingRule
	if changedFiles == 0 && changes == 0 {
		return &defaultRule, defaultLabel, false
	}
	sort.Sort(a.items)
	for _, rule := range a.items {
		if rule.Name == "default" {
			defaultLabel = rule.Label
			defaultExists = true
			defaultRule = rule
		} else if changedFiles <= rule.ChangedFilesThreshold && changes <= rule.ChangesThreshold {
			a.logger.DebugWith(
				log.MetaFields{
					log.F("issue", meta.GetShortName()),
					log.F("rule", rule.Name),
					log.F("file", changedFiles),
					log.F("changed", changes),
				},
				"sizing rule matched")
			return &rule, rule.Label, true
		}
	}
	return &defaultRule, defaultLabel, defaultExists
}

func (s *action) findCurrentMatchedLabel(meta types.Data) (string, bool) {
	for _, label := range meta.GetLabels() {
		if util.StringSliceContains(s.possibleLabels, label) {
			return label, true
		}
	}
	return "", false
}

func (a *action) Apply(state multistep.StateBag) {
	/*
		1. Get number of files and/or changes
		2. Get a list of all the possible applied labels
		3. Check if any update is needed
			3.1 If need different action tag - remove the old one
		4. Update the label
	*/
	a.updatePossibleLabels()
	meta := state.Get("data").(types.Data)
	currentMatchedLabel, exists := a.findCurrentMatchedLabel(meta)
	matchedRule, matchedLabel, matched := a.findMatchedLabel(meta)
	if exists && matched {
		if currentMatchedLabel == matchedLabel {
			a.logger.DebugWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "No need to update label")
			return
		}
		a.logger.DebugWith(
			log.MetaFields{log.F("issue", meta.GetShortName())},
			"Updating label from %s to %s", currentMatchedLabel, matchedLabel)
		meta.RemoveLabel(currentMatchedLabel)
		meta.AddLabel(matchedLabel)
		if matchedRule.Comment != "" {
			meta.AddComment(matchedRule.Comment)
		}
	} else if matched {
		a.logger.DebugWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "Updating label to %s",
			matchedLabel)
		meta.AddLabel(matchedLabel)
		if matchedRule.Comment != "" {
			meta.AddComment(matchedRule.Comment)
		}
	}
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) actions.Action {
	result := action{
		items:  make([]sizingRule, 0),
		logger: log.Get("sizing"),
	}
	for name, internal := range config {
		var item sizingRule
		if e := mapstructure.Decode(internal, &item); e != nil {
			panic(e)
		}
		item.Name = name
		item.Defaults()
		result.items = append(result.items, item)
	}
	return &result
}

func init() {
	actions.RegisterAction("sizing", &factory{})
}
