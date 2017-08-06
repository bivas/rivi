package automerge

import (
	"fmt"
	"strings"

	"github.com/bivas/rivi/engine/actions"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/multistep"
)

type action struct {
	rule *rule
	err  error

	logger log.Logger
}

type MergeableEventData interface {
	Merge(mergeMethod string)
}

type HasReviewersAPIEventData interface {
	GetReviewers() map[string]string
	GetApprovals() []string
}

func (a *action) merge(meta types.EventData) {
	if a.rule.Label == "" {
		mergeable, ok := meta.(MergeableEventData)
		if !ok {
			a.logger.Warning("Event data does not support merge. Check your configurations")
			a.err = fmt.Errorf("Event data does not support merge")
			return
		}
		mergeable.Merge(a.rule.Strategy)
	} else {
		meta.AddLabel(a.rule.Label)
	}
}

func (a *action) getApprovalsFromAPI(meta types.EventData) (int, bool) {
	assigneesList := meta.GetAssignees()
	approvals := 0
	assignees := util.StringSet{}
	assignees.AddAll(assigneesList)
	reviewApi, ok := meta.(HasReviewersAPIEventData)
	if !ok {
		a.logger.WarningWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "Event data does not support reviewers API. Check your configuration")
		a.err = fmt.Errorf("Event data does not support reviewers API")
	} else {
		for _, approver := range reviewApi.GetApprovals() {
			if assignees.Contains(approver) {
				assignees.Remove(approver)
				approvals++
			}
		}
	}
	a.logger.DebugWith(
		log.MetaFields{log.F("issue", meta.GetShortName())},
		"Got %d / %d approvals from API", approvals, len(assigneesList))
	return approvals, assignees.Len() == 0
}

func (a *action) getApprovalsFromComments(meta types.EventData) (int, bool) {
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
	a.logger.DebugWith(
		log.MetaFields{log.F("issue", meta.GetShortName())},
		"Got %d / %d approvals from comments", approvals, len(assigneesList))
	return approvals, assignees.Len() == 0
}

func (a *action) Apply(state multistep.StateBag) {
	meta := state.Get("data").(types.EventData)
	assigneesList := meta.GetAssignees()
	if len(assigneesList) == 0 {
		a.logger.WarningWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "No assignees to issue - skipping")
		return
	}
	calls := []func(types.EventData) (int, bool){
		a.getApprovalsFromAPI,
		a.getApprovalsFromComments,
	}
	for _, call := range calls {
		approvals, all := call(meta)
		if a.rule.Require == 0 && all {
			a.logger.WarningWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "All assignees have approved the PR - merging")
			a.merge(meta)
			return
		} else if a.rule.Require > 0 && approvals >= a.rule.Require {
			a.logger.WarningWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "Got %d required approvals for PR - merging", a.rule.Require)
			a.merge(meta)
			return
		}
	}
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) actions.Action {
	item := rule{}
	if e := mapstructure.Decode(config, &item); e != nil {
		panic(e)
	}
	item.Defaults()
	return &action{rule: &item, logger: log.Get("automerge")}
}

func init() {
	actions.RegisterAction("automerge", &factory{})
}
