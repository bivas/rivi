package slack

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/bivas/rivi/config"
	"github.com/bivas/rivi/engine/actions"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/multistep"
	api "github.com/nlopes/slack"
	"github.com/prometheus/client_golang/prometheus"
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

func (a *action) Apply(state multistep.StateBag) {
	meta := state.Get("data").(types.Data)
	conf := state.Get("config").(config.Configuration)
	if a.rule.Channel == "" && a.rule.Notify == "" {
		a.logger.DebugWith(
			log.MetaFields{log.F("issue", meta.GetShortName())},
			"Skipping action as both Channel and Notify settings are blank")
		return
	}
	if err := a.initClient(conf); err != nil {
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
		a.sendPrivateMessage(conf, meta)
	} else {
		a.sendChannelMessage(conf, meta)
	}
}

func (a *action) postMessage(id string, targets []string, config config.Configuration, meta types.Data) error {
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

func (a *action) sendChannelMessage(config config.Configuration, meta types.Data) {
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
	counter.Inc()
	a.postMessage(channel.ID, meta.GetAssignees(), config, meta)
}

func (a *action) sendPrivateMessage(config config.Configuration, meta types.Data) {
	targets := a.getMessageRecipients(config, meta)
	for _, slacker := range a.toSlackUserId(targets, meta) {
		_, _, id, err := a.client.OpenIMChannel(slacker)
		if err != nil {
			a.logger.WarningWith(log.MetaFields{
				log.F("issue", meta.GetShortName()),
				log.E(err),
				log.F("user.id", slacker)}, "Unable to open IM channel")
			continue
		}
		counter.Inc()
		if err := a.postMessage(id, targets, config, meta); err != nil {
			continue
		}
	}
}

func (a *action) toSlackUserId(users []string, meta types.Data) []string {
	result := make([]string, 0)
	slackers, err := a.client.GetUsers()
	if err != nil {
		a.logger.WarningWith(log.MetaFields{
			log.F("issue", meta.GetShortName()),
			log.E(err)}, "Unable to get users")
		return result
	}
	userSet := util.StringSet{Transformer: strings.ToLower}
	userSet.AddAll(users)
	for _, slacker := range slackers {
		search := strings.ToLower(slacker.Name)
		if userSet.Contains(search) {
			result = append(result, slacker.ID)
		}
	}
	return result
}

func (a *action) getMessageRecipients(config config.Configuration, meta types.Data) []string {
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

func (a *action) initClient(configuration config.Configuration) error {
	inter, err := configuration.GetActionConfig("slack")
	if err != nil {
		return err
	}
	conf, ok := inter.(*actionConfig)
	if !ok {
		return fmt.Errorf("Action Config doesn't match required type (got %T)", inter)
	}
	a.client = api.New(conf.ApiKey)
	a.translator = conf.Translator.Values
	return nil
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) actions.Action {
	item := rule{}
	logger := log.Get("slack")
	if e := mapstructure.Decode(config, &item); e != nil {
		logger.ErrorWith(log.MetaFields{log.E(e)}, "Unable to build action")
		return nil
	}
	item.Defaults()
	return &action{rule: &item, logger: logger}
}

var counter = actions.NewCounter("slack")

func init() {
	actions.RegisterAction("slack", &factory{})
	prometheus.Register(counter)
}
