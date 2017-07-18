package trigger

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
)

type action struct {
	rule   *rule
	client *http.Client
	err    error
	logger log.Logger
}

func (a *action) String() string {
	return fmt.Sprintf("%T{rule: %+v}", *a, a.rule)
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
		a.logger.ErrorWith(
			log.MetaFields{log.F("issue", meta.GetShortName()), log.E(a.err)},
			"Apply got error")
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
		a.logger.ErrorWith(log.MetaFields{log.F("issue", meta.GetShortName()), log.E(e)},
			"Error trying to build trigger request", e)
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
			a.logger.Warning("Skipping header '%s' (non x- headers are not allowed)", name)
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
	return &action{rule: &item, client: http.DefaultClient, logger: log.Get("trigger")}
}

func init() {
	bot.RegisterAction("trigger", &factory{})
}
