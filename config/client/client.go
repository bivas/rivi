package client

import "github.com/spf13/viper"

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

func NewClientConfig(v *viper.Viper) ClientConfig {
	return &clientConfig{v}
}
