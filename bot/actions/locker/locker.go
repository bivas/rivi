package locker

import (
	"github.com/bivas/rivi/bot"
	"github.com/mitchellh/mapstructure"
)

type action struct {
	rule *rule
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	// lock?
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
