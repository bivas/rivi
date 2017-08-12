package slack

import (
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
	"os"
)

type actionConfig struct {
	ApiKey     string `mapstructure:"api-key"`
	Translator struct {
		Values map[string]string
	} `mapstructure:"translator"`
}

func (c *actionConfig) Name() string {
	return "slack"
}

type configBuilder struct {
}

func (c *configBuilder) Build(config map[string]interface{}) (bot.ActionConfig, error) {
	var actionConfig actionConfig
	if err := mapstructure.Decode(config, &actionConfig); err != nil {
		return nil, err
	}
	if env := os.Getenv("RIVI_SLACK_API_KEY"); env != "" {
		log.Get("slack.config").Debug("Setting Slack API-KEY from environment")
		actionConfig.ApiKey = env
	}
	return &actionConfig, nil
}

func init() {
	bot.RegisterActionConfigBuilder("slack", &configBuilder{})
}
