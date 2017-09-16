package client

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ClientConfigTest struct {
	suite.Suite
	token, secret string
}

func (c *ClientConfigTest) TearDownTest() {
	os.Unsetenv("RIVI_CONFIG_TOKEN")
	os.Unsetenv("RIVI_CONFIG_SECRET")
	if c.token != "" {
		os.Setenv("RIVI_CONFIG_TOKEN", c.token)
	}
	if c.secret != "" {
		os.Setenv("RIVI_CONFIG_SECRET", c.secret)
	}
}

func (c *ClientConfigTest) SetupTest() {
	c.token = os.Getenv("RIVI_CONFIG_TOKEN")
	c.secret = os.Getenv("RIVI_CONFIG_SECRET")
	os.Unsetenv("RIVI_CONFIG_TOKEN")
	os.Unsetenv("RIVI_CONFIG_SECRET")
}

func (c *ClientConfigTest) TestClientConfigFromEnv() {
	os.Setenv("RIVI_CONFIG_TOKEN", "token-from-env")
	os.Setenv("RIVI_CONFIG_SECRET", "secret-from-env")
	config := NewClientConfig(viper.New())
	assert.Equal(c.T(), "token-from-env", config.GetOAuthToken(), "token from env")
	assert.Equal(c.T(), "secret-from-env", config.GetSecret(), "secret from env")
}

func (c *ClientConfigTest) TestConfigTest() {
	v := viper.New()
	v.Set("token", "token-from-viper")
	v.Set("secret", "secret-from-viper")
	config := NewClientConfig(v)
	assert.Equal(c.T(), "token-from-viper", config.GetOAuthToken(), "token from viper")
	assert.Equal(c.T(), "secret-from-viper", config.GetSecret(), "secret from viper")
}

func (c *ClientConfigTest) TestConfigFromDevNullFile() {
	config := NewClientConfigFromFile("/dev/null")
	assert.Empty(c.T(), config.GetSecret(), "no secret")
}

func (c *ClientConfigTest) TestConfigFromFileWithConfigPart() {
	config := NewClientConfigFromFile("client_config_test.yml")
	assert.Equal(c.T(), "github-token", config.GetOAuthToken(), "token")
}

func TestClientConfigTest(t *testing.T) {
	suite.Run(t, new(ClientConfigTest))
}
