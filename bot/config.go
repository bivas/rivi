package bot

import (
	"fmt"
	"github.com/bivas/rivi/util"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type ClientConfig interface {
	GetOAuthToken() string
	GetSecret() string
}

type clientConfig struct {
	internal *viper.Viper
}

func (c *clientConfig) GetOAuthToken() string {
	return c.internal.GetString("token")
}

func (c *clientConfig) GetSecret() string {
	return c.internal.GetString("secret")
}

type Configuration interface {
	GetClientConfig() ClientConfig
	GetRoleMembers(role ...string) []string
	GetRoles() []string
	GetRules() []Rule
}

var (
	configSections = []string{"config", "roles", "rules"}
)

type config struct {
	internal     map[string]*viper.Viper
	clientConfig ClientConfig
	rules        []Rule
	roles        map[string][]string
	rolesKeys    []string
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
	set := util.StringSet{}
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
	c.clientConfig = &clientConfig{c.internal["config"]}
	return nil
}

func (c *config) readRolesSection() error {
	c.roles = c.internal["root"].GetStringMapStringSlice("roles")
	util.Logger.Debug("roles from config %s", c.roles)
	c.rolesKeys = make([]string, 0)
	for role := range c.roles {
		c.rolesKeys = append(c.rolesKeys, role)
	}
	return nil
}

func (c *config) readRulesSection() error {
	c.rules = make([]Rule, 0)
	for name := range c.internal["rules"].AllSettings() {
		subname := c.internal["rules"].Sub(name)
		r := &rule{
			name:      name,
			condition: buildConditionFromConfiguration(subname),
			action:    buildActionFromConfiguration(subname),
		}
		util.Logger.Debug("appending rule %s", r)
		c.rules = append(c.rules, r)
	}
	return nil
}

func (c *config) readSections() error {
	if err := c.readConfigSection(); err != nil {
		return err
	}
	if err := c.readRolesSection(); err != nil {
		return err
	}
	if err := c.readRulesSection(); err != nil {
		return err
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
