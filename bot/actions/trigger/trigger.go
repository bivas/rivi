package trigger

import (
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type action struct {
	rule *rule
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	message := &message{
		Time:   time.Now(),
		Number: meta.GetNumber(),
		Title:  meta.GetTitle(),
		State:  meta.GetState(),
		Owner:  meta.GetOwner(),
		Repo:   meta.GetRepo(),
		Origin: meta.GetOrigin(),
	}
	util.Logger.Debug("Prepared a message %+v", message)
	request := &http.Request{

	}
	http.DefaultClient.Do(request)

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
	bot.RegisterAction("trigger", &factory{})
}
