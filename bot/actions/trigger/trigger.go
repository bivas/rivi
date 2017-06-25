package trigger

import (
	"fmt"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strings"
	"time"
)

type action struct {
	rule   *rule
	client *http.Client
	err    error
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	request := a.prepareRequest(meta)
	response, e := a.client.Do(request)
	if e != nil {
		a.err = fmt.Errorf("Triggering to %s, resulted in error. %s",
			a.rule.Endpoint,
			e)
	} else if response.StatusCode >= 400 {
		a.err = fmt.Errorf("Triggering to %s, resulted in wrong status code. %d",
			a.rule.Endpoint,
			response.StatusCode)
	}

	if a.err != nil {
		util.Logger.Error("%s", a.err)
	}

}
func (a *action) prepareRequest(meta bot.EventData) *http.Request {
	message := &message{
		Time:   time.Now(),
		Number: meta.GetNumber(),
		Title:  meta.GetTitle(),
		State:  meta.GetState(),
		Owner:  meta.GetOwner(),
		Repo:   meta.GetRepo(),
		Origin: meta.GetOrigin(),
	}
	body := processMessage(&a.rule.Body, message)
	request, e := http.NewRequest(a.rule.Method, a.rule.Endpoint, body)
	if e != nil {
		util.Logger.Error("Error trying to build trigger request. %s", e)
	}
	a.setRequestHeaders(request)
	return request
}
func (a *action) setRequestHeaders(request *http.Request) {
	request.Header.Set("User-Agent", "RiviBot-Agent/1.0")
	request.Header.Set("X-RiviBot-Event", "trigger")
	request.Header.Set("Content-Type", "application/json")
	for name, value := range a.rule.Headers {
		lowerName := strings.ToLower(name)
		if !strings.HasPrefix(lowerName, "x-") || strings.HasPrefix(lowerName, "x-rivibot") {
			util.Logger.Warning("Skipping header '%s' (non x- headers are not allowed)", name)
		} else {
			request.Header.Set(name, value)
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
	return &action{rule: &item, client: http.DefaultClient}
}

func init() {
	bot.RegisterAction("trigger", &factory{})
}
