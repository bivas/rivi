package automerge

import (
	"fmt"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type action struct {
	rule *rule
	err  error
}

type MergeableEventData interface {
	Merge(mergeMethod string)
}

func (a *action) merge(meta bot.EventData) {
	if a.rule.Label == "" {
		mergeable, ok := meta.(MergeableEventData)
		if !ok {
			util.Logger.Warning("Event data does not support merge. Check your configurations")
			a.err = fmt.Errorf("Event data does not support merge")
			return
		}
		mergeable.Merge(a.rule.Strategy)
	} else {
		meta.AddLabel(a.rule.Label)
	}
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	assigneesList := meta.GetAssignees()
	if len(assigneesList) {
		util.Logger.Debug("No assignees to issue - skipping")
		return
	}
	approvals := 0
	assignees := util.StringSet{}
	assignees.AddAll(assigneesList)
	for _, comment := range meta.GetComments() {
		if !assignees.Contains(comment.Commenter) {
			continue
		}
		clean := strings.ToLower(strings.TrimSpace(comment.Comment))
		if _, ok := approvedSearchPhrases[clean]; ok {
			assignees.Remove(comment.Commenter)
			approvals++
		}
	}
	if a.rule.Require == 0 && assignees.Len() == 0 {
		util.Logger.Debug("All assignees have approved the PR - merging")
		a.merge(meta)
	} else if a.rule.Require > 0 && approvals >= a.rule.Require {
		util.Logger.Debug("Got %d required approvals for PR - merging", a.rule.Require)
		a.merge(meta)
	}
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) bot.Action {
	item := rule{}
	if e := mapstructure.Decode(config, &item); e != nil {
		panic(e)
	}
	item.Defaults()
	return &action{rule: &item}
}

func init() {
	bot.RegisterAction("automerge", &factory{})
}
