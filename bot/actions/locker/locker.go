package locker

import (
	"fmt"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"github.com/mitchellh/mapstructure"
)

type LockableEventData interface {
	Lock()
	Unlock()
	LockState() bool
}

type action struct {
	rule *rule
	err  error
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	lockable, ok := meta.(LockableEventData)
	if !ok {
		util.Logger.Warning("Event data does not support locking. Check your configurations")
		a.err = fmt.Errorf("Event data does not support locking")
		return
	}
	if lockable.LockState() {
		util.Logger.Debug("Issue is locked")
		if a.rule.State == "unlock" || a.rule.State == "change" {
			util.Logger.Debug("unlocking issue %d", meta.GetNumber())
			lockable.Unlock()
		} else if a.rule.State == "lock" {
			util.Logger.Debug("Issue %d is already locked - nothing changed", meta.GetNumber())
		}
	} else {
		util.Logger.Debug("Issue is unlocked")
		if a.rule.State == "lock" || a.rule.State == "change" {
			util.Logger.Debug("Locking issue %d", meta.GetNumber())
			lockable.Lock()
		} else if a.rule.State == "lock" {
			util.Logger.Debug("Issue %d is already unlocked - nothing changed", meta.GetNumber())
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
	return &action{rule: &item}
}

func init() {
	bot.RegisterAction("locker", &factory{})
}
