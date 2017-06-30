package sizing

import (
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"github.com/mitchellh/mapstructure"
	"sort"
)

type action struct {
	items          rules
	possibleLabels []string
}

func (s *action) updatePossibleLabels() {
	set := util.StringSet{}
	for _, item := range s.items {
		set.Add(item.Label)
	}
	s.possibleLabels = set.Values()
}

func (s *action) findMatchedLabel(meta bot.EventData) (*sizingRule, string, bool) {
	changedFiles := meta.GetChangedFiles()
	add, del := meta.GetChanges()
	changes := add + del
	defaultLabel := ""
	defaultExists := false
	var defaultRule sizingRule
	sort.Sort(s.items)
	for _, rule := range s.items {
		if rule.Name == "default" {
			defaultLabel = rule.Label
			defaultExists = true
			defaultRule = rule
		} else if changedFiles <= rule.ChangedFilesThreshold && changes <= rule.ChangesThreshold {
			util.Logger.Debug("[action] [(%d) %s] sizing rule %s matched with %d files and %d changes",
				meta.GetNumber(),
				meta.GetTitle(),
				rule.Name,
				changedFiles,
				changes)
			return &rule, rule.Label, true
		}
	}
	return &defaultRule, defaultLabel, defaultExists
}

func (s *action) findCurrentMatchedLabel(meta bot.EventData) (string, bool) {
	for _, label := range meta.GetLabels() {
		if util.StringSliceContains(s.possibleLabels, label) {
			return label, true
		}
	}
	return "", false
}

func (s *action) Apply(config bot.Configuration, meta bot.EventData) {
	/*
		1. Get number of files and/or changes
		2. Get a list of all the possible applied labels
		3. Check if any update is needed
			3.1 If need different action tag - remove the old one
		4. Update the label
	*/
	s.updatePossibleLabels()
	currentMatchedLabel, exists := s.findCurrentMatchedLabel(meta)
	matchedRule, matchedLabel, matched := s.findMatchedLabel(meta)
	if exists && matched {
		if currentMatchedLabel == matchedLabel {
			util.Logger.Debug("[action] [(%d) %s] No need to update label",
				meta.GetNumber(),
				meta.GetTitle())
			return
		}
		util.Logger.Debug("[action] [(%d) %s] Updating label from %s to %s",
			meta.GetNumber(),
			meta.GetTitle(),
			currentMatchedLabel,
			matchedLabel)
		meta.RemoveLabel(currentMatchedLabel)
		meta.AddLabel(matchedLabel)
		if matchedRule.Comment != "" {
			meta.AddComment(matchedRule.Comment)
		}
	} else if matched {
		util.Logger.Debug("[action] [(%d) %s] Updating label to %s",
			meta.GetNumber(),
			meta.GetTitle(),
			matchedLabel)
		meta.AddLabel(matchedLabel)
		if matchedRule.Comment != "" {
			meta.AddComment(matchedRule.Comment)
		}
	}
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) bot.Action {
	result := action{
		items: make([]sizingRule, 0),
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
	bot.RegisterAction("sizing", &factory{})
}
