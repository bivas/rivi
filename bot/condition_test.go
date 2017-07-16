package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockConditionEventData struct {
	Number         int
	Owner          string
	Repo           string
	Labels         []string
	FileNames      []string
	FileExtensions []string
	Title          string
	Description    string
	Ref            string
	RawPayload     []byte
}

func (m *mockConditionEventData) GetRawPayload() []byte {
	return m.RawPayload
}

func (m *mockConditionEventData) GetRef() string {
	return m.Ref
}

func (m *mockConditionEventData) GetNumber() int {
	return m.Number
}

func (m *mockConditionEventData) GetTitle() string {
	return m.Title
}

func (m *mockConditionEventData) GetDescription() string {
	return m.Description
}

func (m *mockConditionEventData) GetState() string {
	panic("implement me")
}

func (*mockConditionEventData) GetOrigin() string {
	panic("implement me")
}

func (m *mockConditionEventData) GetOwner() string {
	return m.Owner
}

func (m *mockConditionEventData) GetRepo() string {
	return m.Repo
}

func (m *mockConditionEventData) GetLabels() []string {
	return m.Labels
}

func (*mockConditionEventData) HasLabel(label string) bool {
	panic("implement me")
}

func (*mockConditionEventData) AddLabel(label string) {
	panic("implement me")
}

func (*mockConditionEventData) RemoveLabel(label string) {
	panic("implement me")
}

func (*mockConditionEventData) GetAssignees() []string {
	panic("implement me")
}

func (*mockConditionEventData) HasAssignee(assignee string) bool {
	panic("implement me")
}

func (*mockConditionEventData) AddAssignees(assignees ...string) {
	panic("implement me")
}

func (*mockConditionEventData) RemoveAssignees(assignees ...string) {
	panic("implement me")
}

func (*mockConditionEventData) GetComments() []Comment {
	panic("implement me")
}

func (*mockConditionEventData) AddComment(comment string) {
	panic("implement me")
}

func (m *mockConditionEventData) GetFileNames() []string {
	return m.FileNames
}

func (*mockConditionEventData) GetChangedFiles() int {
	panic("implement me")
}

func (m *mockConditionEventData) GetFileExtensions() []string {
	return m.FileExtensions
}

func (*mockConditionEventData) GetChanges() (int, int) {
	panic("implement me")
}

func getConfig(t *testing.T) Configuration {
	c, err := newConfiguration("config_test.yml")
	if err != nil {
		t.Fatalf("Got error during config read. %s", err)
	}
	return c
}

func TestMatchLabel(t *testing.T) {
	c := getConfig(t)
	meta := &mockConditionEventData{Labels: []string{"label1", "pending-approval"}}
	matched := make([]string, 0)
	for _, rule := range c.GetRules() {
		if rule.Accept(meta) {
			matched = append(matched, rule.Name())
		}
	}
	assert.Contains(t, matched, "rule1", "matched")
	assert.NotContains(t, matched, "rule2", "matched")
}

func TestSkipLabel(t *testing.T) {
	c := getConfig(t)
	meta := &mockConditionEventData{Labels: []string{"pending-approval"}}
	matched := make([]string, 0)
	for _, rule := range c.GetRules() {
		if rule.Accept(meta) {
			matched = append(matched, rule.Name())
		}
	}
	assert.Len(t, matched, 0, "matched")
}

func TestMatchPattern(t *testing.T) {
	c := getConfig(t)
	meta := &mockConditionEventData{
		Labels: []string{"pending-approval"},
		FileNames: []string{
			"foo.txt",
			"path/to/docs/foo.txt",
		}}
	matched := make([]string, 0)
	for _, rule := range c.GetRules() {
		if rule.Accept(meta) {
			matched = append(matched, rule.Name())
		}
	}
	assert.Len(t, matched, 1, "matched")
	assert.Contains(t, matched, "rule4", "matched")
}

func TestMatchExt(t *testing.T) {
	c := getConfig(t)
	meta := &mockConditionEventData{FileExtensions: []string{".scala", ".go"}}
	matched := false
	for _, rule := range c.GetRules() {
		if rule.Name() == "rule3" {
			matched = true
			assert.True(t, rule.Accept(meta), "extension")
		}
	}
	assert.True(t, matched, "matched")
}

func TestTitleStartsWith(t *testing.T) {
	var rule rule
	rule.condition.Title.StartsWith = "BUGFIX"
	meta := &mockConditionEventData{Title: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Title = "BUGFIX it"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestTitleEndsWith(t *testing.T) {
	var rule rule
	rule.condition.Title.EndsWith = "WIP"
	meta := &mockConditionEventData{Title: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Title = "BUGFIX WIP"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestTitlePattern(t *testing.T) {
	var rule rule
	rule.condition.Title.Patterns = []string{".* Bug( )?[0-9]{5} .*"}
	meta := &mockConditionEventData{Title: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Title = "This PR for Bug1 with comment"
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Title = "This PR for Bug 45456 with comment"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestDescriptionStartsWith(t *testing.T) {
	var rule rule
	rule.condition.Description.StartsWith = "BUGFIX"
	meta := &mockConditionEventData{Description: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Description = "BUGFIX it"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestDescriptionEndsWith(t *testing.T) {
	var rule rule
	rule.condition.Description.EndsWith = "WIP"
	meta := &mockConditionEventData{Description: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Description = "BUGFIX WIP"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestDescriptionPattern(t *testing.T) {
	var rule rule
	rule.condition.Description.Patterns = []string{"(?s)~~~.*deps:"}
	meta := &mockConditionEventData{Description: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Description = "~~~\n     test_priorities"
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Description = "~~~\n    deps:\nplenty of dependencies"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestMatchEmptyCondition(t *testing.T) {
	meta := &mockConditionEventData{}
	rule := rule{condition: Condition{}}
	assert.True(t, rule.Accept(meta), "none")
}

func TestMatchRef(t *testing.T) {
	var rule rule
	rule.condition.Ref.Equals = "master"
	meta := &mockConditionEventData{Ref: "development"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Ref = "master"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestRefPatters(t *testing.T) {
	var rule rule
	rule.condition.Ref.Patterns = []string{"integration-v[0-9]{2}$"}
	meta := &mockConditionEventData{Ref: "development"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Ref = "integration-v11"
	assert.True(t, rule.Accept(meta), "should match")
}
