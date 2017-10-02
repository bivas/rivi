package client

import (
	"github.com/bivas/rivi/util/log"
	"github.com/spf13/viper"
)

type ClientConfig interface {
	GetOAuthToken() string
	GetSecret() string
	GetApplicationID() int
	GetPrivateKeyFile() string
}

type clientConfig struct {
	*viper.Viper
}

func (c *clientConfig) GetPrivateKeyFile() string {
	return c.GetString("private_key_file")
}

func (c *clientConfig) GetApplicationID() int {
	return c.GetInt("appid")
}

func (c *clientConfig) GetOAuthToken() string {
	return c.GetString("token")
}

func (c *clientConfig) GetSecret() string {
	return c.GetString("secret")
}

func NewDefaultClientConfig() ClientConfig {
	return NewClientConfig(viper.New())
}

func NewClientConfig(v *viper.Viper) ClientConfig {
	v.SetEnvPrefix("rivi_config")
	v.BindEnv("token")
	v.BindEnv("secret")
	v.BindEnv("appid")
	v.BindEnv("private_key_file")
	return &clientConfig{v}
}

func NewClientConfigFromFile(file string) ClientConfig {
	logger := log.Get("config.client")
	v := viper.New()
	logger.DebugWith(
		log.MetaFields{
			log.F("config", file),
		}, "Loading client config from file",
	)
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		logger.WarningWith(
			log.MetaFields{
				log.E(err),
				log.F("config", file),
			}, "Error loading config from file",
		)
		return NewDefaultClientConfig()
	}
	configViper := v.Sub("config")
	if configViper != nil {
		logger.DebugWith(
			log.MetaFields{
				log.F("config", file),
			}, "Loading client config from sub config",
		)
		return NewClientConfig(configViper)
	}
	logger.DebugWith(
		log.MetaFields{
			log.F("config", file),
		}, "Loading client config",
	)
	return NewClientConfig(v)
}
