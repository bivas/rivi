package bot

import (
	"github.com/stretchr/testify/assert"
	. "testing"
)

func assertClientConfig(t *T, config ClientConfig) {
	assert.Equal(t, config.GetOAuthToken(), "github-token", "oath token")
	assert.Equal(t, config.GetSecret(), "github-secret", "secret")
}

func assertRoles(t *T, configuration Configuration) {
	roles := configuration.GetRoles()
	assert.Contains(t, roles, "admins", "roles")
	assert.Contains(t, roles, "reviewers", "roles")
	assert.Contains(t, roles, "testers", "roles")
	assert.NotContains(t, roles, "dummy", "dummy role")

	admins := configuration.GetRoleMembers("admins")
	assert.Contains(t, admins, "user1", "admin.user1")
	assert.Contains(t, admins, "user2", "admin.user2")
	assert.NotContains(t, admins, "user3", "admin.user3")

	reviewers := configuration.GetRoleMembers("reviewers")
	assert.Contains(t, reviewers, "user1", "reviewers.user1")
	assert.NotContains(t, reviewers, "user2", "reviewers.user2")
	assert.Contains(t, reviewers, "user3", "reviewers.user3")

	testers := configuration.GetRoleMembers("testers")
	assert.Contains(t, testers, "user2", "testers.user2")
	assert.NotContains(t, testers, "user1", "reviewers.user3")
	assert.NotContains(t, testers, "user3", "reviewers.user3")
}

func assertRules(t *T, configuration Configuration) {
	rules := configuration.GetRules()
	assert.Len(t, rules, 4, "rules len")
	ruleNames := make([]string, 0)
	for _, rule := range rules {
		ruleNames = append(ruleNames, rule.Name())
	}
	assert.Contains(t, ruleNames, "rule1", "rule name")
	assert.Contains(t, ruleNames, "rule2", "rule name")
	assert.Contains(t, ruleNames, "rule3", "rule name")
	assert.Contains(t, ruleNames, "rule4", "rule name")
}

func TestReadConfig(t *T) {
	c, err := newConfiguration("config_test.yml")
	if err != nil {
		t.Fatalf("Got error during config read. %s", err)
	}
	assertClientConfig(t, c.GetClientConfig())
	assertRoles(t, c)
	assertRules(t, c)
}
