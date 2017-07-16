package automerge

import (
	"fmt"
	"strings"

	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"github.com/mitchellh/mapstructure"
)

type action struct {
	rule *rule
	err  error
}

type MergeableEventData interface {
	Merge(mergeMethod string)
}

type HasReviewersAPIEventData interface {
	GetReviewers() map[string]string
	GetApprovals() []string
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

func (a *action) getApprovalsFromAPI(meta bot.EventData) (int, bool) {
	assigneesList := meta.GetAssignees()
	approvals := 0
	assignees := util.StringSet{}
	assignees.AddAll(assigneesList)
	reviewApi, ok := meta.(HasReviewersAPIEventData)
	if !ok {
		util.Logger.Warning("Event data does not support reviewers API. Check your configuration")
		a.err = fmt.Errorf("Event data does not support reviewers API")
	} else {
		for _, approver := range reviewApi.GetApprovals() {
			if assignees.Contains(approver) {
				assignees.Remove(approver)
				approvals++
			}
		}
	}
	return approvals, assignees.Len() == 0
}

func (a *action) getApprovalsFromComments(meta bot.EventData) (int, bool) {
	assigneesList := meta.GetAssignees()
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
	return approvals, assignees.Len() == 0
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	assigneesList := meta.GetAssignees()
	if len(assigneesList) == 0 {
		util.Logger.Debug("No assignees to issue - skipping")
		return
	}
	calls := []func(bot.EventData) (int, bool){
		a.getApprovalsFromAPI,
		a.getApprovalsFromComments,
	}
	for _, call := range calls {
		approvals, all := call(meta)
		if a.rule.Require == 0 && all {
			util.Logger.Debug("All assignees have approved the PR - merging")
			a.merge(meta)
			return
		} else if a.rule.Require > 0 && approvals >= a.rule.Require {
			util.Logger.Debug("Got %d required approvals for PR - merging", a.rule.Require)
			a.merge(meta)
			return
		}
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
