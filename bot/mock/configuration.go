package mock

import (
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
)

type MockClientConfig struct {
	OAuthToken string
	Secret     string
}

func (m *MockClientConfig) GetOAuthToken() string {
	return m.GetOAuthToken()
}

func (m *MockClientConfig) GetSecret() string {
	return m.Secret
}

type MockConfiguration struct {
	MockClientConfig *MockClientConfig
	RoleMembers      map[string][]string
}

func (m *MockConfiguration) GetClientConfig() bot.ClientConfig {
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

func (m *MockConfiguration) GetRules() []bot.Rule {
	panic("implement me")
}
