package autoassign

import (
	"math/rand"
	"strings"

	"fmt"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
)

type action struct {
	rule *rule

	logger log.Logger
}

func (a *action) String() string {
	return fmt.Sprintf("%T{rule: %+v}", *a, a.rule)
}

func (a *action) findAssignedRoles(assignees []string, config bot.Configuration) []string {
	var assignedRoles []string
	if len(assignees) > 0 {
		assignedRolesSet := util.StringSet{}
		for _, role := range config.GetRoles() {
			for _, member := range config.GetRoleMembers(role) {
				for _, assignee := range assignees {
					if member == assignee {
						assignedRolesSet.Add(role)
					}
				}
			}
		}
		assignedRoles = assignedRolesSet.Values()
		a.logger.Debug("There are %d assignees from %d roles", len(assignees), len(assignedRoles))
	}
	return assignedRoles
}

func (a *action) findLookupRoles(config bot.Configuration, assignedRoles []string) []string {
	lookupRoles := config.GetRoles()
	if len(a.rule.FromRoles) > 0 {
		lookupRoles = a.rule.FromRoles
	}
	a.logger.Debug("lookup roles are %s", lookupRoles)
	return lookupRoles
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	assignees := meta.GetAssignees()
	if len(assignees) >= a.rule.Require {
		a.logger.DebugWith(
			log.MetaFields{log.F("issue", meta.GetShortName())},
			"Skipping as there are matched required assignees")
		return
	}
	assignedRoles := a.findAssignedRoles(assignees, config)
	lookupRoles := a.findLookupRoles(config, assignedRoles)

	winners := a.randomUsers(config, meta, lookupRoles)
	if len(winners) > 0 {
		meta.AddAssignees(winners...)
	}
}

func (a *action) randomUsers(config bot.Configuration, meta bot.EventData, lookupRoles []string) []string {
	possibleSet := util.StringSet{Transformer: strings.ToLower}
	possibleSet.AddAll(config.GetRoleMembers(lookupRoles...)).Remove(meta.GetOrigin())
	for _, assignee := range meta.GetAssignees() {
		possibleSet.Remove(assignee)
	}
	a.logger.DebugWith(
		log.MetaFields{log.F("issue", meta.GetShortName())},
		"There are %d possible assignees from %d roles", possibleSet.Len(), len(lookupRoles))
	if possibleSet.Len() == 0 {
		return []string{}
	}
	remainingRequired := a.rule.Require - len(meta.GetAssignees())
	if remainingRequired < 0 {
		remainingRequired = 0
	}
	a.logger.DebugWith(
		log.MetaFields{log.F("issue", meta.GetShortName())},
		"Require %d assignees", remainingRequired)
	possible := possibleSet.Values()
	winners := make([]string, remainingRequired)
	for i := 0; i < remainingRequired; i++ {
		index := rand.Intn(len(possible))
		if possible[index] == "" {
			i--
		} else {
			winners[i] = possible[index]
			possible[index] = ""
		}
	}
	a.logger.DebugWith(
		log.MetaFields{log.F("issue", meta.GetShortName())},
		"Selecting users %s as assignees", winners)
	return winners
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) bot.Action {
	item := rule{}
	if e := mapstructure.Decode(config, &item); e != nil {
		panic(e)
	}
	item.Defaults()
	return &action{rule: &item, logger: log.Get("autoassign")}
}

func init() {
	bot.RegisterAction("autoassign", &factory{})
}
