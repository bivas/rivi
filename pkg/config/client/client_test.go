package client

import (
	"os"
	"testing"

	"github.com/spf13/viper"
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
	c.Require().Equal("token-from-env", config.GetOAuthToken(), "token from env")
	c.Require().Equal("secret-from-env", config.GetSecret(), "secret from env")
}

func (c *ClientConfigTest) TestConfigTest() {
	v := viper.New()
	v.Set("token", "token-from-viper")
	v.Set("secret", "secret-from-viper")
	config := NewClientConfig(v)
	c.Require().Equal("token-from-viper", config.GetOAuthToken(), "token from viper")
	c.Require().Equal("secret-from-viper", config.GetSecret(), "secret from viper")
}

func (c *ClientConfigTest) TestConfigFromDevNullFile() {
	config := NewClientConfigFromFile("/dev/null")
	c.Require().Empty(config.GetSecret(), "no secret")
}

func (c *ClientConfigTest) TestConfigFromFileWithConfigPart() {
	config := NewClientConfigFromFile("client_config_test.yml")
	c.Require().Equal("github-token", config.GetOAuthToken(), "token")
}

func TestClientConfigTest(t *testing.T) {
	suite.Run(t, new(ClientConfigTest))
}
