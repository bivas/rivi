package automerge

import (
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/bot/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockMergableEventData struct {
	mock.MockEventData
	merged bool
	method string
}

func (m *mockMergableEventData) Merge(mergeMethod string) {
	m.merged = true
	m.method = mergeMethod
}

func TestNotCapableToMerge(t *testing.T) {
	action := action{rule: &rule{}}
	action.rule.Defaults()
	config := &mock.MockConfiguration{}
	meta := &mock.MockEventData{Assignees: []string{"user1"}, Comments: []bot.Comment{
		{Commenter: "user1", Comment: "approved"},
	}}
	action.Apply(config, meta)
	assert.NotNil(t, action.err, "should be unable to merge")
}

func TestShouldNotMergeMergeMissingApprovals(t *testing.T) {
	action := action{rule: &rule{}}
	action.rule.Defaults()
	config := &mock.MockConfiguration{}
	mockEventData := mock.MockEventData{Assignees: []string{"user1", "user2"}, Comments: []bot.Comment{
		{Commenter: "user1", Comment: "approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(config, meta)
	assert.False(t, meta.merged, "should not be merged")
}

func TestCapableToMerge(t *testing.T) {
	action := action{rule: &rule{}}
	action.rule.Defaults()
	config := &mock.MockConfiguration{}
	mockEventData := mock.MockEventData{Assignees: []string{"user1"}, Comments: []bot.Comment{
		{Commenter: "user1", Comment: "approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(config, meta)
	assert.True(t, meta.merged, "should be merged")
	assert.Equal(t, "merge", meta.method, "default should be merge")
}

func TestOriginComment(t *testing.T) {
	action := action{rule: &rule{}}
	action.rule.Defaults()
	config := &mock.MockConfiguration{}
	mockEventData := mock.MockEventData{Assignees: []string{"user1"}, Origin: "user2", Comments: []bot.Comment{
		{Commenter: "user2", Comment: "approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(config, meta)
	assert.False(t, meta.merged, "should not be merged")
}

func TestNotApprovedComment(t *testing.T) {
	action := action{rule: &rule{}}
	action.rule.Defaults()
	config := &mock.MockConfiguration{}
	mockEventData := mock.MockEventData{Assignees: []string{"user1"}, Origin: "user2", Comments: []bot.Comment{
		{Commenter: "user1", Comment: "not approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(config, meta)
	assert.False(t, meta.merged, "should not be merged")
}

func TestSameApprovedComment(t *testing.T) {
	action := action{rule: &rule{Require: 2}}
	action.rule.Defaults()
	config := &mock.MockConfiguration{}
	mockEventData := mock.MockEventData{Assignees: []string{"user1"}, Origin: "user2", Comments: []bot.Comment{
		{Commenter: "user1", Comment: "approved"},
		{Commenter: "user1", Comment: "approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(config, meta)
	assert.False(t, meta.merged, "should not be merged")
}

func TestLabel(t *testing.T) {
	action := action{rule: &rule{Label: "approved"}}
	action.rule.Defaults()
	config := &mock.MockConfiguration{}
	mockEventData := mock.MockEventData{Assignees: []string{"user1"}, Origin: "user2", Comments: []bot.Comment{
		{Commenter: "user1", Comment: "approved"},
	}}
	meta := &mockMergableEventData{MockEventData: mockEventData}
	action.Apply(config, meta)
	assert.False(t, meta.merged, "should not be merged")
	assert.Len(t, meta.AddedLabels, 1, "should label and not merge")
	assert.Equal(t, "approved", meta.AddedLabels[0], "approved")
}
