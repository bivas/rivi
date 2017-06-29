package autoassign

import (
	"github.com/bivas/rivi/bot/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActionApplyRequire1NoAssignees(t *testing.T) {
	action := action{rule: &rule{Require: 1}}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2", "user3"}
	config := &mock.MockConfiguration{RoleMembers: roles}
	meta := &mock.MockEventData{Assignees: []string{}}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedAssignees, 1, "assignment")
}

func TestActionApplyRequire2With1Assignee(t *testing.T) {
	action := action{rule: &rule{Require: 2}}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2", "user3"}
	config := &mock.MockConfiguration{RoleMembers: roles}
	meta := &mock.MockEventData{Assignees: []string{"user1"}}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedAssignees, 1, "assignment")
	assert.Len(t, meta.Assignees, 2, "required")
	assert.Contains(t, meta.Assignees, "user1", "original")
}

func TestActionApplyRequireFromRole(t *testing.T) {
	action := action{rule: &rule{Require: 1, FromRoles: []string{"group"}}}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2", "user3"}
	roles["group"] = []string{"user4"}
	config := &mock.MockConfiguration{RoleMembers: roles}
	meta := &mock.MockEventData{Assignees: []string{}}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedAssignees, 1, "assignment")
	assert.NotContains(t, meta.AddedAssignees, "user1", "matched group")
	assert.NotContains(t, meta.AddedAssignees, "user2", "matched group")
	assert.NotContains(t, meta.AddedAssignees, "user3", "matched group")
	assert.Contains(t, meta.AddedAssignees, "user4", "matched group")
}

func TestActionApplyWithoutOrigin(t *testing.T) {
	action := action{rule: &rule{Require: 1, FromRoles: []string{"default"}}}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2"}
	for i := 0; i < 10; i++ {
		config := &mock.MockConfiguration{RoleMembers: roles}
		meta := &mock.MockEventData{Assignees: []string{}, Origin: "user1"}
		action.Apply(config, meta)
		assert.Len(t, meta.AddedAssignees, 1, "assignment")
		assert.NotContains(t, meta.AddedAssignees, "user1", "matched group")
		assert.Contains(t, meta.AddedAssignees, "user2", "matched group")
	}
}

func TestActionApplyHasAssignee(t *testing.T) {
	action := action{rule: &rule{Require: 1, FromRoles: []string{"group"}}}
	roles := make(map[string][]string)
	roles["default"] = []string{"user1", "user2", "user3"}
	roles["group"] = []string{"user4"}
	config := &mock.MockConfiguration{RoleMembers: roles}
	meta := &mock.MockEventData{Assignees: []string{"user1"}}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedAssignees, 0, "assignment")
}
