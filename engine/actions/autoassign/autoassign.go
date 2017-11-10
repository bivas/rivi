package autoassign

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bivas/rivi/config"
	"github.com/bivas/rivi/engine/actions"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/multistep"
	"github.com/prometheus/client_golang/prometheus"
)

type action struct {
	rule *rule

	logger log.Logger
}

func (a *action) String() string {
	return fmt.Sprintf("%T{rule: %+v}", *a, a.rule)
}

func (a *action) findAssignedRoles(assignees []string, config config.Configuration) []string {
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

func (a *action) findLookupRoles(config config.Configuration, assignedRoles []string) []string {
	lookupRoles := config.GetRoles()
	if len(a.rule.FromRoles) > 0 {
		lookupRoles = a.rule.FromRoles
	}
	a.logger.Debug("lookup roles are %s", lookupRoles)
	return lookupRoles
}

func (a *action) Apply(state multistep.StateBag) {
	meta := state.Get("data").(types.Data)
	conf := state.Get("config").(config.Configuration)
	assignees := meta.GetAssignees()
	if len(assignees) >= a.rule.Require {
		a.logger.DebugWith(
			log.MetaFields{log.F("issue", meta.GetShortName())},
			"Skipping as there are matched required assignees")
		return
	}
	assignedRoles := a.findAssignedRoles(assignees, conf)
	lookupRoles := a.findLookupRoles(conf, assignedRoles)

	winners := a.randomUsers(conf, meta, lookupRoles)
	if len(winners) > 0 {
		counter.Inc()
		meta.AddAssignees(winners...)
	}
}

func (a *action) randomUsers(config config.Configuration, meta types.Data, lookupRoles []string) []string {
	possibleSet := util.StringSet{Transformer: strings.ToLower}
	possibleSet.AddAll(config.GetRoleMembers(lookupRoles...)).Remove(meta.GetOrigin().User)
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

func (*factory) BuildAction(config map[string]interface{}) actions.Action {
	item := rule{}
	logger := log.Get("autoassign")
	if e := mapstructure.Decode(config, &item); e != nil {
		logger.ErrorWith(log.MetaFields{log.E(e)}, "Unable to build action")
		return nil
	}
	item.Defaults()
	return &action{rule: &item, logger: logger}
}

var counter = actions.NewCounter("autoassign")

func init() {
	actions.RegisterAction("autoassign", &factory{})
	prometheus.Register(counter)
}
