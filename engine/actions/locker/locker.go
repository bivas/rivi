package locker

import (
	"fmt"

	"github.com/bivas/rivi/engine/actions"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/multistep"
)

type LockableData interface {
	Lock()
	Unlock()
	LockState() bool
}

type action struct {
	rule   *rule
	err    error
	logger log.Logger
}

func (a *action) Apply(state multistep.StateBag) {
	meta := state.Get("data").(types.Data)
	lockable, ok := meta.(LockableData)
	if !ok {
		a.logger.WarningWith(
			log.MetaFields{log.F("issue", meta.GetShortName())},
			"Event data does not support locking. Check your configurations")
		a.err = fmt.Errorf("Event data does not support locking")
		return
	}
	if lockable.LockState() {
		a.logger.Debug("Issue is locked")
		if a.rule.State == "unlock" || a.rule.State == "change" {
			a.logger.DebugWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "unlocking issue")
			lockable.Unlock()
		} else if a.rule.State == "lock" {
			a.logger.DebugWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "Issue is already locked - nothing changed")
		}
	} else {
		a.logger.Debug("Issue is unlocked")
		if a.rule.State == "lock" || a.rule.State == "change" {
			a.logger.DebugWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "Locking issue")
			lockable.Lock()
		} else if a.rule.State == "lock" {
			a.logger.DebugWith(log.MetaFields{log.F("issue", meta.GetShortName())}, "Issue is already unlocked - nothing changed")
		}
	}
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) actions.Action {
	item := rule{}
	logger := log.Get("locker")
	if e := mapstructure.Decode(config, &item); e != nil {
		logger.ErrorWith(log.MetaFields{log.E(e)}, "Unable to build action")
		return nil
	}
	return &action{rule: &item, logger: logger}
}

func init() {
	actions.RegisterAction("locker", &factory{})
}
