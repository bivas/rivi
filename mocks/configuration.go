package mocks

import (
	"github.com/bivas/rivi/config/action"
	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/engine"
	"github.com/bivas/rivi/util"
)

type MockClientConfig struct {
	OAuthToken     string
	Secret         string
	ApplicationID  int
	PrivateKeyFile string
}

func (m *MockClientConfig) GetPrivateKeyFile() string {
	return m.PrivateKeyFile
}

func (m *MockClientConfig) GetApplicationID() int {
	return m.ApplicationID
}

func (m *MockClientConfig) GetOAuthToken() string {
	return m.OAuthToken
}

func (m *MockClientConfig) GetSecret() string {
	return m.Secret
}

type MockConfiguration struct {
	MockClientConfig *MockClientConfig
	RoleMembers      map[string][]string
}

func (m *MockConfiguration) GetActionConfig(kind string) (action.ActionConfig, error) {
	panic("implement me")
}

func (m *MockConfiguration) GetClientConfig() client.ClientConfig {
	return m.MockClientConfig
}

func (m *MockConfiguration) GetRoleMembers(roles ...string) []string {
	result := make([]string, 0)
	for _, role := range roles {
		if members, ok := m.RoleMembers[role]; ok {
			result = append(result, members...)
		}
	}
	set := util.StringSet{}
	set.AddAll(result)
	return set.Values()
}

func (m *MockConfiguration) GetRoles() []string {
	result := make([]string, 0)
	for role := range m.RoleMembers {
		result = append(result, role)
	}
	return result
}

func (m *MockConfiguration) GetRules() []engine.Rule {
	panic("implement me")
}
