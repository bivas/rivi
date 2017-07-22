package slack

import (
	"fmt"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
	api "github.com/nlopes/slack"
	"text/template"
)

type action struct {
	rule       *rule
	client     *api.Client
	translator map[string]string
	template   *template.Template
	logger     log.Logger
}

func (a *action) String() string {
	return fmt.Sprintf("%T{rule: %+v}", *a, a.rule)
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	if a.rule.Channel == "" && a.rule.Notify == "" {
		a.logger.DebugWith(
			log.MetaFields{log.F("issue", meta.GetShortName())},
			"Skipping action as both Channel and Notify settings are blank")
		return
	}
	if err := a.initClient(config); err != nil {
		a.logger.WarningWith(
			log.MetaFields{log.F("issue", meta.GetShortName()), log.E(err)},
			"Unable to init Slack client")
		return
	}
	if err := a.compileMessage(); err != nil {
		a.logger.WarningWith(
			log.MetaFields{log.F("issue", meta.GetShortName()), log.E(err)},
			"Unable to compile Slack message")
		return
	}
	if a.rule.Channel == "" {
		a.sendPrivateMessage(config, meta)
	} else {
		a.sendChannelMessage(config, meta)
	}
}

func (a *action) postMessage(id string, targets []string, config bot.Configuration, meta bot.EventData) error {
	text, err := serializeMessage(a.template, buildMessage(meta, targets))
	if err != nil {
		a.logger.WarningWith(
			log.MetaFields{log.F("issue", meta.GetShortName()), log.E(err)},
			"Unable to build text message")
		return err
	}
	if _, resp, err := a.client.PostMessage(id, text, api.NewPostMessageParameters()); err != nil {
		a.logger.WarningWith(
			log.MetaFields{log.F("issue", meta.GetShortName()), log.E(err)},
			"Unable to post message")
		a.logger.DebugWith(
			log.MetaFields{log.F("issue", meta.GetShortName()), log.F("response", resp)},
			"Unable to post")
		return err
	}
	return nil
}

func (a *action) sendChannelMessage(config bot.Configuration, meta bot.EventData) {
	channel, err := a.client.GetChannelInfo(a.rule.Channel)
	if err != nil {
		a.logger.WarningWith(
			log.MetaFields{
				log.F("issue", meta.GetShortName()),
				log.E(err),
				log.F("channel", a.rule.Channel)},
			"Unable to get channel info")
		return
	}
	a.postMessage(channel.ID, meta.GetAssignees(), config, meta)
}

func (a *action) sendPrivateMessage(config bot.Configuration, meta bot.EventData) {
	targets := a.getMessageRecipients(config, meta)
	for _, slacker := range targets {
		_, _, id, err := a.client.OpenIMChannel(slacker)
		if err != nil {
			a.logger.WarningWith(log.MetaFields{
				log.F("issue", meta.GetShortName()),
				log.E(err),
				log.F("user", slacker)}, "Unable to open IM channel")
			continue
		}
		if err := a.postMessage(id, targets, config, meta); err != nil {
			continue
		}
	}
}

func (a *action) getMessageRecipients(config bot.Configuration, meta bot.EventData) []string {
	var result []string
	if a.rule.Notify == "assignees" {
		result = meta.GetAssignees()
	} else {
		result = config.GetRoleMembers(a.rule.Notify)
	}
	if len(a.translator) > 0 {
		for index, user := range result {
			if slacker, exists := a.translator[user]; exists {
				result[index] = slacker
			}
		}
	}
	return result
}


func (a *action) compileMessage() error {
	t, err := template.New("slack-action").Parse(a.rule.Message)
	if err != nil {
		return err
	}
	a.template = t
	return nil
}

func (a *action) initClient(configuration bot.Configuration) error {
	inter, err := configuration.GetActionConfig("slack")
	if err != nil {
		return err
	}
	config, ok := inter.(*actionConfig)
	if !ok {
		return fmt.Errorf("Action Config doesn't match required type (got %T)", inter)
	}
	a.client = api.New(config.ApiKey)
	a.translator = config.Translator.Values
	return nil
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) bot.Action {
	item := rule{}
	if e := mapstructure.Decode(config, &item); e != nil {
		panic(e)
	}
	item.Defaults()
	return &action{rule: &item, logger: log.Get("slack")}
}

func init() {
	bot.RegisterAction("slack", &factory{})
}
