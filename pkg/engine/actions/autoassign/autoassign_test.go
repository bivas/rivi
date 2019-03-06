package autoassign

import (
	"testing"

	"github.com/bivas/rivi/pkg/mocks"
	"github.com/bivas/rivi/pkg/types"
	"github.com/bivas/rivi/pkg/util/log"
	"github.com/bivas/rivi/pkg/util/state"
	"github.com/stretchr/testify/assert"
)

func TestSerialization(t *testing.T) {
	input := map[string]interface{}{
		"roles":   []string{"role1", "role2"},
		"require": 2,
	}

	var f factory
	result := f.BuildAction(input)
	assert.NotNil(t, result, "should create action")
	s, ok := result.(*action)
	assert.True(t, ok, "should be of this package")
	assert.Equal(t, 2, s.rule.Require, "require")
	assert.EqualValues(t, []string{"role1", "role2"}, s.rule.FromRoles, "roles")
}

func TestActionApplyRequire1NoAssignees(t *testing.T) {
	action := action{rule: &rule{Require: 1}, logger: log.Get("autoassign.test")}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2", "user3"}
	config := &mocks.MockConfiguration{RoleMembers: roles}
	meta := &mocks.MockData{Assignees: []string{}}
	action.Apply(state.New(config, meta))
	assert.Len(t, meta.AddedAssignees, 1, "assignment")
}

func TestActionApplyRequire2With1Assignee(t *testing.T) {
	action := action{rule: &rule{Require: 2}, logger: log.Get("autoassign.test")}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2", "user3"}
	config := &mocks.MockConfiguration{RoleMembers: roles}
	meta := &mocks.MockData{Assignees: []string{"user1"}}
	action.Apply(state.New(config, meta))
	assert.Len(t, meta.AddedAssignees, 1, "assignment")
	assert.Len(t, meta.Assignees, 2, "required")
	assert.Contains(t, meta.Assignees, "user1", "original")
}

func TestActionApplyRequireFromRole(t *testing.T) {
	action := action{rule: &rule{Require: 1, FromRoles: []string{"group"}}, logger: log.Get("autoassign.test")}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2", "user3"}
	roles["group"] = []string{"user4"}
	config := &mocks.MockConfiguration{RoleMembers: roles}
	meta := &mocks.MockData{Assignees: []string{}}
	action.Apply(state.New(config, meta))
	assert.Len(t, meta.AddedAssignees, 1, "assignment")
	assert.NotContains(t, meta.AddedAssignees, "user1", "matched group")
	assert.NotContains(t, meta.AddedAssignees, "user2", "matched group")
	assert.NotContains(t, meta.AddedAssignees, "user3", "matched group")
	assert.Contains(t, meta.AddedAssignees, "user4", "matched group")
}

func TestActionApplyWithoutOrigin(t *testing.T) {
	action := action{rule: &rule{Require: 1, FromRoles: []string{"default"}}, logger: log.Get("autoassign.test")}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2"}
	for i := 0; i < 10; i++ {
		config := &mocks.MockConfiguration{RoleMembers: roles}
		meta := &mocks.MockData{Assignees: []string{}, Origin: types.Origin{User: "user1"}}
		action.Apply(state.New(config, meta))
		assert.Len(t, meta.AddedAssignees, 1, "assignment")
		assert.NotContains(t, meta.AddedAssignees, "user1", "matched group")
		assert.Contains(t, meta.AddedAssignees, "user2", "matched group")
	}
}

func TestActionApplyWithoutOriginCaseInsensitive(t *testing.T) {
	action := action{rule: &rule{Require: 1, FromRoles: []string{"default"}}, logger: log.Get("autoassign.test")}
	roles := make(map[string][]string)
	roles["default"] = []string{"usEr1", "user2"}
	for i := 0; i < 10; i++ {
		config := &mocks.MockConfiguration{RoleMembers: roles}
		meta := &mocks.MockData{Assignees: []string{}, Origin: types.Origin{User: "user1"}}
		action.Apply(state.New(config, meta))
		assert.Len(t, meta.AddedAssignees, 1, "assignment")
		assert.NotContains(t, meta.AddedAssignees, "user1", "matched group")
		assert.Contains(t, meta.AddedAssignees, "user2", "matched group")
	}
}

func TestActionApplyHasAssignee(t *testing.T) {
	action := action{rule: &rule{Require: 1, FromRoles: []string{"group"}}, logger: log.Get("autoassign.test")}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2", "user3"}
	roles["group"] = []string{"user4"}
	config := &mocks.MockConfiguration{RoleMembers: roles}
	meta := &mocks.MockData{Assignees: []string{"user1"}}
	action.Apply(state.New(config, meta))
	assert.Len(t, meta.AddedAssignees, 0, "assignment")
}

func TestActionApplyAllAssignees(t *testing.T) {
	action := action{rule: &rule{Require: 1, FromRoles: []string{"group"}}, logger: log.Get("autoassign.test")}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2"}
	config := &mocks.MockConfiguration{RoleMembers: roles}
	meta := &mocks.MockData{Assignees: []string{"user1", "user2"}}
	action.Apply(state.New(config, meta))
	assert.Len(t, meta.AddedAssignees, 0, "assignment")
}
