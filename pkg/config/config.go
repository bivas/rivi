package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bivas/rivi/pkg/config/action"
	"github.com/bivas/rivi/pkg/config/client"
	"github.com/bivas/rivi/pkg/engine"
	"github.com/bivas/rivi/pkg/util"
	"github.com/bivas/rivi/pkg/util/log"
	"github.com/spf13/viper"
)

type Configuration interface {
	GetClientConfig() client.ClientConfig
	GetRoleMembers(role ...string) []string
	GetRoles() []string
	GetRules() []engine.Rule
	GetActionConfig(kind string) (action.ActionConfig, error)
}

var (
	configSections = []string{"config", "roles", "rules"}
	lc             = log.Get("config")
)

type config struct {
	internal      map[string]*viper.Viper
	clientConfig  client.ClientConfig
	rules         []engine.Rule
	roles         map[string][]string
	rolesKeys     []string
	actionConfigs map[string]action.ActionConfig
}

func (c *config) GetActionConfig(kind string) (action.ActionConfig, error) {
	config, exists := c.actionConfigs[kind]
	if !exists {
		return nil, fmt.Errorf("No such action config %s", kind)
	}
	return config, nil
}

func (c *config) GetClientConfig() client.ClientConfig {
	return c.clientConfig
}

func (c *config) GetRoleMembers(roles ...string) []string {
	result := make([]string, 0)
	for _, role := range roles {
		if members, ok := c.roles[role]; ok {
			result = append(result, members...)
		}
	}
	set := util.StringSet{Transformer: strings.ToLower}
	set.AddAll(result)
	return set.Values()
}

func (c *config) GetRoles() []string {
	return c.rolesKeys
}

func (c *config) GetRules() []engine.Rule {
	return c.rules
}

func (c *config) readConfigSection() error {
	internal := c.internal["config"]
	if internal == nil {
		internal = viper.New()
	}
	c.clientConfig = client.NewClientConfig(internal)
	return nil
}

func (c *config) readRolesSection() error {
	c.roles = c.internal["root"].GetStringMapStringSlice("roles")
	lc.Debug("roles from config %s", c.roles)
	c.rolesKeys = make([]string, 0)
	for role := range c.roles {
		c.rolesKeys = append(c.rolesKeys, role)
	}
	log.DebugWith(log.MetaFields{log.F("roles", c.rolesKeys)}, "Loaded %d roles", len(c.rolesKeys))
	return nil
}

func (c *config) readRulesSection() error {
	c.rules = make([]engine.Rule, 0)
	if c.internal["rules"] == nil {
		return nil
	}
	for name := range c.internal["rules"].AllSettings() {
		subname := c.internal["rules"].Sub(name)
		r := engine.NewRule(name, subname)
		lc.Debug("appending rule %s", r)
		c.rules = append(c.rules, r)
	}
	sort.Sort(engine.RulesByConditionOrder(c.rules))
	lc.Debug("Loaded %d rules", len(c.rules))
	return nil
}

func (c *config) readSections() error {
	sections := []func() error{
		c.readConfigSection,
		c.readRolesSection,
		c.readRulesSection,
	}
	for _, section := range sections {
		if err := section(); err != nil {
			lc.DebugWith(log.MetaFields{log.F("section", section), log.E(err)}, "Section got an error")
			return err
		}
	}
	c.actionConfigs = action.BuildActionConfigs(c.internal["root"])
	return nil
}

func (c *config) readConfiguration(configPath string) error {
	c.internal["root"] = viper.New()
	c.internal["root"].SetConfigName("rivi")
	c.internal["root"].SetConfigFile(configPath)

	if err := c.internal["root"].ReadInConfig(); err != nil {
		return err
	}
	for _, section := range configSections {
		c.internal[section] = c.internal["root"].Sub(section)
	}
	return c.readSections()
}

func NewConfiguration(configPath string) (Configuration, error) {
	c := &config{
		internal: map[string]*viper.Viper{},
	}
	if err := c.readConfiguration(configPath); err != nil {
		return nil, err
	}
	return c, nil
}
