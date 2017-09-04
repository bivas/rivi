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
	internal *viper.Viper
}

func (c *clientConfig) GetPrivateKeyFile() string {
	return c.internal.GetString("private_key_file")
}

func (c *clientConfig) GetApplicationID() int {
	return c.internal.GetInt("appid")
}

func (c *clientConfig) GetOAuthToken() string {
	return c.internal.GetString("token")
}

func (c *clientConfig) GetSecret() string {
	return c.internal.GetString("secret")
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
		return NewClientConfig(viper.New())
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
