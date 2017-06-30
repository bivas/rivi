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

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	mergeable, ok := meta.(MergeableEventData)
	if !ok {
		util.Logger.Warning("Event data does not support merge. Check your configurations")
		a.err = fmt.Errorf("Event data does not support merge")
		return
	}
	approvals := 0
	assignees := util.StringSet{}
	assignees.AddAll(meta.GetAssignees())
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
	if approvals >= a.rule.Require {
		mergeable.Merge(a.rule.Strategy)
	}
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) bot.Action {
	item := rule{}
	if e := mapstructure.Decode(config, &item); e != nil {
		panic(e)
	}
	return &action{rule: &item}
}

func init() {
	bot.RegisterAction("automerge", &factory{})
}
