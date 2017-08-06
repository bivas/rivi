package config

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bivas/rivi/config/action"
	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/engine"
	"github.com/bivas/rivi/util"
	"github.com/bivas/rivi/util/log"
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

func (c *config) getSection(path string) (string, *viper.Viper) {
	split := strings.SplitN(path, ".", 2)
	item, exists := c.internal[split[0]]
	if exists {
		return split[1], item
	} else {
		return path, c.internal["root"]
	}
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
	lc.Debug("Loaded %d roles", len(c.rolesKeys))
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
			return err
		}
	}
	c.actionConfigs = action.BuildActionConfigs(c.internal["root"])
	return nil
}

func (c *config) readConfiguration(configPath string) error {
	c.internal["root"] = viper.New()
	c.internal["root"].SetConfigName("bot")
	c.internal["root"].SetConfigFile(configPath)

	if err := c.internal["root"].ReadInConfig(); err != nil {
		return err
	}
	rootConfigFir := filepath.Dir(configPath)
	for _, section := range configSections {
		sectionInclude := c.internal["root"].GetString(fmt.Sprintf("%s.include", section))
		if sectionInclude != "" {
			lc.Debug("Attempt loading %s config from file %s", section, sectionInclude)
			c.internal[section] = viper.New()
			c.internal[section].SetConfigFile(filepath.Join(rootConfigFir, sectionInclude))
			if err := c.internal[section].ReadInConfig(); err != nil {
				return err
			}
		} else {
			c.internal[section] = c.internal["root"].Sub(section)
		}
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
