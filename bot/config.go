package bot

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bivas/rivi/util"
	"github.com/bivas/rivi/util/log"
	"github.com/spf13/viper"
)

type ClientConfig interface {
	GetOAuthToken() string
	GetSecret() string
}

type clientConfig struct {
	internal *viper.Viper
}

func (c *clientConfig) GetOAuthToken() string {
	c.internal.SetEnvPrefix("rivi_config")
	c.internal.BindEnv("token")
	return c.internal.GetString("token")
}

func (c *clientConfig) GetSecret() string {
	c.internal.SetEnvPrefix("rivi_config")
	c.internal.BindEnv("secret")
	return c.internal.GetString("secret")
}

type ActionConfig interface {
	Name() string
}

type ActionConfigBuilder interface {
	Build(config map[string]interface{}) (ActionConfig, error)
}

type Configuration interface {
	GetClientConfig() ClientConfig
	GetRoleMembers(role ...string) []string
	GetRoles() []string
	GetRules() []Rule
	GetActionConfig(kind string) (ActionConfig, error)
}

var (
	configSections       = []string{"config", "roles", "rules"}
	actionConfigBuilders = make(map[string]ActionConfigBuilder)
)

func RegisterActionConfigBuilder(name string, builder ActionConfigBuilder) {
	search := strings.ToLower(name)
	_, exists := actionConfigBuilders[search]
	if exists {
		log.Error("action config build for %s exists", name)
	} else {
		log.Debug("registering action config builder %s", name)
		actionConfigBuilders[search] = builder
	}
}

type config struct {
	internal      map[string]*viper.Viper
	clientConfig  ClientConfig
	rules         []Rule
	roles         map[string][]string
	rolesKeys     []string
	actionConfigs map[string]ActionConfig
}

func (c *config) GetActionConfig(kind string) (ActionConfig, error) {
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

func (c *config) GetClientConfig() ClientConfig {
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

func (c *config) GetRules() []Rule {
	return c.rules
}

func (c *config) readConfigSection() error {
	internal := c.internal["config"]
	if internal == nil {
		internal = viper.New()
	}
	c.clientConfig = &clientConfig{internal}
	return nil
}

func (c *config) readRolesSection() error {
	c.roles = c.internal["root"].GetStringMapStringSlice("roles")
	log.Debug("roles from config %s", c.roles)
	c.rolesKeys = make([]string, 0)
	for role := range c.roles {
		c.rolesKeys = append(c.rolesKeys, role)
	}
	log.Debug("Loaded %d roles", len(c.rolesKeys))
	return nil
}

func (c *config) readRulesSection() error {
	c.rules = make([]Rule, 0)
	for name := range c.internal["rules"].AllSettings() {
		subname := c.internal["rules"].Sub(name)
		r := &rule{
			name:      name,
			condition: buildConditionFromConfiguration(subname),
			actions:   buildActionsFromConfiguration(subname),
		}
		log.Debug("appending rule %s", r)
		c.rules = append(c.rules, r)
	}
	sort.Sort(rulesByConditionOrder(c.rules))
	log.Debug("Loaded %d rules", len(c.rules))
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
	c.actionConfigs = make(map[string]ActionConfig)
	for kind, builder := range actionConfigBuilders {
		log.Debug("Building configuration for %s", kind)
		actionConfig := c.internal["root"].GetStringMap(kind)
		if actionConfig == nil || len(actionConfig) == 0 {
			log.Warning("No matching section for %s", kind)
			continue
		}
		config, err := builder.Build(actionConfig)
		if err != nil {
			log.ErrorWith(log.MetaFields{log.E(err)}, "Error while building %s config", kind)
			continue
		}
		c.actionConfigs[kind] = config
	}
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
			log.Debug("Attempt loading %s config from file %s", section, sectionInclude)
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

func newConfiguration(configPath string) (Configuration, error) {
	c := &config{
		internal: map[string]*viper.Viper{},
	}
	if err := c.readConfiguration(configPath); err != nil {
		return nil, err
	}
	return c, nil
}
