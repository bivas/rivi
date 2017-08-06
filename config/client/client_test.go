package client

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClientConfigFromEnv(t *testing.T) {
	os.Setenv("RIVI_CONFIG_TOKEN", "token-from-env")
	os.Setenv("RIVI_CONFIG_SECRET", "secret-from-env")
	config := NewClientConfig(viper.New())
	assert.Equal(t, "token-from-env", config.GetOAuthToken(), "token from env")
	assert.Equal(t, "secret-from-env", config.GetSecret(), "secret from env")
}

func TestConfigTest(t *testing.T) {
	os.Unsetenv("RIVI_CONFIG_TOKEN")
	os.Unsetenv("RIVI_CONFIG_SECRET")
	v := viper.New()
	v.Set("token", "token-from-viper")
	v.Set("secret", "secret-from-viper")
	config := NewClientConfig(v)
	assert.Equal(t, "token-from-viper", config.GetOAuthToken(), "token from viper")
	assert.Equal(t, "secret-from-viper", config.GetSecret(), "secret from viper")
}
