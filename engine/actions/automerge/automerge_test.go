package automerge

import (
	"testing"

	"github.com/bivas/rivi/mocks"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/bivas/rivi/util/state"
	"github.com/stretchr/testify/assert"
)

type mockMergableEventData struct {
	mocks.MockEventData
	merged    bool
	method    string
	reviewers map[string]string
	approvals []string
}

func (m *mockMergableEventData) GetReviewers() map[string]string {
	return m.reviewers
}

func (m *mockMergableEventData) GetApprovals() []string {
	return m.approvals
}

func (m *mockMergableEventData) Merge(mergeMethod string) {
	m.merged = true
	m.method = mergeMethod
}

func TestSerialization(t *testing.T) {
	input := map[string]interface{}{
		"strategy": "squash",
		"require":  2,
		"label":    "label1",
	}

	var f factory
	result := f.BuildAction(input)
	assert.NotNil(t, result, "should create action")
	s, ok := result.(*action)
	assert.True(t, ok, "should be of this package")
	assert.Equal(t, 2, s.rule.Require, "require")
	assert.Equal(t, "squash", s.rule.Strategy, "strategy")
	assert.Equal(t, "label1", s.rule.Label, "label")
}

func TestSerializationDefaults(t *testing.T) {
	input := map[string]interface{}{
		"require": 2,
		"label":   "label1",
	}

	var f factory
	result := f.BuildAction(input)
	assert.NotNil(t, result, "should create action")
	s, ok := result.(*action)
	assert.True(t, ok, "should be of this package")
	assert.Equal(t, 2, s.rule.Require, "require")
	assert.Equal(t, "merge", s.rule.Strategy, "strategy")
	assert.Equal(t, "label1", s.rule.Label, "label")
}

func TestNoReviewersAPI(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	meta := &mocks.MockEventData{Assignees: []string{"user1"}, Comments: []types.Comment{
		{Commenter: "user1", Comment: "approved"},
	}}
	action.Apply(state.New(config, meta))
	assert.NotNil(t, action.err, "should be error on api")
}

func TestMergeWithAPI(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	mockEventData := mocks.MockEventData{Assignees: []string{"user1"}}
	meta := &mockMergableEventData{MockEventData: mockEventData, approvals: []string{"user1"}}
	action.Apply(state.New(config, meta))
	assert.True(t, meta.merged, "should be merged")
	assert.Equal(t, "merge", meta.method, "default should be merge")
}

func TestMergeWithAPINoApprovals(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	mockEventData := mocks.MockEventData{Assignees: []string{"user1"}}
	meta := &mockMergableEventData{MockEventData: mockEventData, approvals: []string{"user2"}}
	action.Apply(state.New(config, meta))
	assert.False(t, meta.merged, "should be merged")
}

func TestNotCapableToMerge(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	meta := &mocks.MockEventData{Assignees: []string{"user1"}, Comments: []types.Comment{
		{Commenter: "user1", Comment: "approved"},
	}}
	action.Apply(state.New(config, meta))
	assert.NotNil(t, action.err, "should be unable to merge")
}

func TestShouldNotMergeMergeMissingApprovals(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	mockEventData := mocks.MockEventData{Assignees: []string{"user1", "user2"}, Comments: []types.Comment{
		{Commenter: "user1", Comment: "approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(state.New(config, meta))
	assert.False(t, meta.merged, "should not be merged")
}

func TestCapableToMerge(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	mockEventData := mocks.MockEventData{Assignees: []string{"user1"}, Comments: []types.Comment{
		{Commenter: "user1", Comment: "approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(state.New(config, meta))
	assert.True(t, meta.merged, "should be merged")
	assert.Equal(t, "merge", meta.method, "default should be merge")
}

func TestOriginComment(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	mockEventData := mocks.MockEventData{Assignees: []string{"user1"}, Origin: "user2", Comments: []types.Comment{
		{Commenter: "user2", Comment: "approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(state.New(config, meta))
	assert.False(t, meta.merged, "should not be merged")
}

func TestNotApprovedComment(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	mockEventData := mocks.MockEventData{Assignees: []string{"user1"}, Origin: "user2", Comments: []types.Comment{
		{Commenter: "user1", Comment: "not approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(state.New(config, meta))
	assert.False(t, meta.merged, "should not be merged")
}

func TestSameApprovedComment(t *testing.T) {
	action := action{rule: &rule{Require: 2}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	mockEventData := mocks.MockEventData{Assignees: []string{"user1"}, Origin: "user2", Comments: []types.Comment{
		{Commenter: "user1", Comment: "approved"},
		{Commenter: "user1", Comment: "approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(state.New(config, meta))
	assert.False(t, meta.merged, "should not be merged")
}

func TestLabel(t *testing.T) {
	action := action{rule: &rule{Label: "approved"}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	mockEventData := mocks.MockEventData{Assignees: []string{"user1"}, Origin: "user2", Comments: []types.Comment{
		{Commenter: "user1", Comment: "approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(state.New(config, meta))
	assert.False(t, meta.merged, "should not be merged")
	assert.Len(t, meta.AddedLabels, 1, "should label and not merge")
	assert.Equal(t, "approved", meta.AddedLabels[0], "approved")
}

func TestNoAssignees(t *testing.T) {
	action := action{rule: &rule{}, logger: log.Get("automerge.test")}
	action.rule.Defaults()
	config := &mocks.MockConfiguration{}
	mockEventData := mocks.MockEventData{Assignees: []string{}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(state.New(config, meta))
	assert.False(t, meta.merged, "should not be merged")
}
