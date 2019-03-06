package action

import (
	"strings"

	"github.com/bivas/rivi/pkg/util/log"
	"github.com/spf13/viper"
)

type ActionConfig interface {
	Name() string
}

type ActionConfigBuilder interface {
	Build(config map[string]interface{}) (ActionConfig, error)
}

var (
	actionConfigBuilders = make(map[string]ActionConfigBuilder)
	lca                  = log.Get("config.action")
)

func RegisterActionConfigBuilder(name string, builder ActionConfigBuilder) {
	search := strings.ToLower(name)
	_, exists := actionConfigBuilders[search]
	if exists {
		lca.Error("action config build for %s exists", name)
	} else {
		lca.Debug("registering action config builder %s", name)
		actionConfigBuilders[search] = builder
	}
}

func BuildActionConfigs(config *viper.Viper) map[string]ActionConfig {
	actionConfigs := make(map[string]ActionConfig)
	for kind, builder := range actionConfigBuilders {
		lca.Debug("Building configuration for %s", kind)
		actionConfig := config.GetStringMap(kind)
		if actionConfig == nil || len(actionConfig) == 0 {
			lca.Warning("No matching section for %s", kind)
			continue
		}
		config, err := builder.Build(actionConfig)
		if err != nil {
			lca.ErrorWith(log.MetaFields{log.E(err)}, "Error while building %s config", kind)
			continue
		}
		actionConfigs[kind] = config
	}
	return actionConfigs
}
