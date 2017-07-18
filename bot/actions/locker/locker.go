package locker

import (
	"fmt"

	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
)

type LockableEventData interface {
	Lock()
	Unlock()
	LockState() bool
}

type action struct {
	rule   *rule
	err    error
	logger log.Logger
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	lockable, ok := meta.(LockableEventData)
	if !ok {
		a.logger.WarningWith(
			log.MetaFields{{"issue", meta.GetShortName()}},
			"Event data does not support locking. Check your configurations")
		a.err = fmt.Errorf("Event data does not support locking")
		return
	}
	if lockable.LockState() {
		a.logger.Debug("Issue is locked")
		if a.rule.State == "unlock" || a.rule.State == "change" {
			a.logger.DebugWith(log.MetaFields{{"issue", meta.GetShortName()}}, "unlocking issue")
			lockable.Unlock()
		} else if a.rule.State == "lock" {
			a.logger.DebugWith(log.MetaFields{{"issue", meta.GetShortName()}}, "Issue is already locked - nothing changed")
		}
	} else {
		a.logger.Debug("Issue is unlocked")
		if a.rule.State == "lock" || a.rule.State == "change" {
			a.logger.DebugWith(log.MetaFields{{"issue", meta.GetShortName()}}, "Locking issue")
			lockable.Lock()
		} else if a.rule.State == "lock" {
			a.logger.DebugWith(log.MetaFields{{"issue", meta.GetShortName()}}, "Issue is already unlocked - nothing changed")
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
	return &action{rule: &item, logger: log.Get("locker")}
}

func init() {
	bot.RegisterAction("locker", &factory{})
}
