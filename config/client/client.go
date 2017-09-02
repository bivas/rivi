package client

import "github.com/spf13/viper"

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
