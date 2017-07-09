package slack

import (
	"fmt"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"github.com/mitchellh/mapstructure"
	api "github.com/nlopes/slack"
	"text/template"
)

type action struct {
	rule       *rule
	client     *api.Client
	translator map[string]string
	template   *template.Template
}

func (a *action) String() string {
	return fmt.Sprintf("%T{rule: %+v}", *a, a.rule)
}

func (a *action) Apply(config bot.Configuration, meta bot.EventData) {
	if a.rule.Channel == "" && a.rule.Notify == "" {
		util.Logger.Debug("Skipping action as both Channel and Notify settings are blank")
		return
	}
	if err := a.initClient(config); err != nil {
		util.Logger.Warning("Unable to init Slack client. %s", err)
		return
	}
	if err := a.compileMessage(); err != nil {
		util.Logger.Warning("Unable to init Slack message. %s", err)
		return
	}
	if a.rule.Channel == "" {
		a.sendPrivateMessage(config, meta)
	} else {
		a.sendChannelMessage(config, meta)
	}
}

func (a *action) postMessage(id, slacker string, config bot.Configuration, meta bot.EventData) error {
	message := buildFromEventData(meta, slacker)
	text, err := buildMessage(a.template, message)
	if err != nil {
		util.Logger.Warning("Unable to build text message for user '%s'. %s", slacker, err)
		return err
	}
	if _, resp, err := a.client.PostMessage(id, text, api.NewPostMessageParameters()); err != nil {
		util.Logger.Warning("Unable to post message to '%s'. %s", slacker, err)
		util.Logger.Debug("Unable to post to '%s'. %s", slacker, resp)
		return err
	}
	return nil
}

func (a *action) sendChannelMessage(config bot.Configuration, meta bot.EventData) {
	channel, err := a.client.GetChannelInfo(a.rule.Channel)
	if err != nil {
		util.Logger.Warning("Unable to get channel '%s' info. %s", a.rule.Channel, err)
		return
	}
	a.postMessage(channel.ID, "", config, meta)
}

func (a *action) sendPrivateMessage(config bot.Configuration, meta bot.EventData) {
	targets := a.getMessageRecipients(config, meta)
	for _, slacker := range targets {
		_, _, id, err := a.client.OpenIMChannel(slacker)
		if err != nil {
			util.Logger.Warning("Unable to open IM channel to '%s'. %s", slacker, err)
			continue
		}
		if err := a.postMessage(id, slacker, config, meta); err != nil {
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
	return &action{rule: &item}
}

func init() {
	bot.RegisterAction("slack", &factory{})
}
