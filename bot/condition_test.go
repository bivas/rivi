package bot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockConditionEventData struct {
	Labels         []string
	FileNames      []string
	FileExtensions []string
	Title          string
}

func (*mockConditionEventData) GetNumber() int {
	panic("implement me")
}

func (m *mockConditionEventData) GetTitle() string {
	return m.Title
}

func (m *mockConditionEventData) GetState() string {
	panic("implement me")
}

func (*mockConditionEventData) GetOrigin() string {
	panic("implement me")
}

func (*mockConditionEventData) GetOwner() string {
	panic("implement me")
}

func (*mockConditionEventData) GetRepo() string {
	panic("implement me")
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

func TestMatchEmptyCondition(t *testing.T) {
	meta := &mockConditionEventData{}
	rule := rule{condition: Condition{}}
	assert.True(t, rule.Accept(meta), "none")
}
