package sizing

import (
	"testing"

	"github.com/bivas/rivi/bot/mock"
	"github.com/bivas/rivi/util/log"
	"github.com/stretchr/testify/assert"
)

func buildRule(name, label string, changedFiles int) *sizingRule {
	result := &sizingRule{Name: name, Label: label, ChangedFilesThreshold: changedFiles, Comment: name}
	result.Defaults()
	return result
}

func buildRules(withDefault bool) *action {
	result := &action{items: rules{
		*buildRule("l", "size/l", 150),
		*buildRule("xs", "size/xs", 5),
		*buildRule("s", "size/s", 15),
		*buildRule("m", "size/m", 75),
	}, possibleLabels: []string{
		"size/xs",
		"size/s",
		"size/m",
	}, logger: log.Get("sizing.test")}
	if withDefault {
		result.possibleLabels = append(result.possibleLabels, "default-label")
		result.items = append(result.items, sizingRule{Name: "default", Label: "default-label"})
	}
	return result
}

func TestSerialization(t *testing.T) {
	input := map[string]interface{}{
		"size1": map[string]interface{}{
			"label":                   "label1",
			"comment":                 "comment1",
			"changed-files-threshold": 10,
		},
		"size2": map[string]interface{}{
			"label":             "label2",
			"comment":           "comment2",
			"changes-threshold": 20,
		},
	}

	var f factory
	result := f.BuildAction(input)
	assert.NotNil(t, result, "should create action")
	s, ok := result.(*action)
	assert.True(t, ok, "should be of this package")
	assert.Equal(t, 2, len(s.items), "rules")
	if s.items[0].Name == "size1" {
		assert.Equal(t, "label1", s.items[0].Label, "label")
		assert.Equal(t, "comment1", s.items[0].Comment, "comment")
		assert.Equal(t, 10, s.items[0].ChangedFilesThreshold, "file")
		assert.Equal(t, "label2", s.items[1].Label, "label")
		assert.Equal(t, "comment2", s.items[1].Comment, "comment")
		assert.Equal(t, 20, s.items[1].ChangesThreshold, "changes")
	}
	if s.items[0].Name == "size2" {
		assert.Equal(t, "label2", s.items[0].Label, "label")
		assert.Equal(t, "comment2", s.items[0].Comment, "comment")
		assert.Equal(t, 20, s.items[0].ChangesThreshold, "changes")
		assert.Equal(t, "label1", s.items[1].Label, "label")
		assert.Equal(t, "comment1", s.items[1].Comment, "comment")
		assert.Equal(t, 10, s.items[1].ChangedFilesThreshold, "file")
	}

}

func TestSizingXS(t *testing.T) {
	action := buildRules(false)
	config := &mock.MockConfiguration{}
	meta := &mock.MockEventData{Labels: []string{}, ChangedFiles: 1, ChangesAdd: 24, ChangesRemove: 31}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "labels")
	assert.Equal(t, "size/xs", meta.AddedLabels[0], "size label")
	assert.Len(t, meta.AddedComments, 1, "comments")
	assert.Equal(t, "xs", meta.AddedComments[0], "comment")
}

func TestSizing(t *testing.T) {
	action := buildRules(false)
	config := &mock.MockConfiguration{}
	meta := &mock.MockEventData{Labels: []string{}, ChangedFiles: 8}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "labels")
	assert.Equal(t, "size/s", meta.AddedLabels[0], "size label")
}

func TestWithDefaultShouldNotApply(t *testing.T) {
	action := buildRules(true)
	config := &mock.MockConfiguration{}
	meta := &mock.MockEventData{Labels: []string{}, ChangedFiles: 8}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "labels")
	assert.Equal(t, "size/s", meta.AddedLabels[0], "size label")
}

func TestWithDefaultShouldApply(t *testing.T) {
	action := buildRules(true)
	config := &mock.MockConfiguration{}
	meta := &mock.MockEventData{Labels: []string{}, ChangedFiles: 800}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "labels")
	assert.Equal(t, "default-label", meta.AddedLabels[0], "size label")
}

func TestSizingUpdateNeeded(t *testing.T) {
	action := buildRules(false)
	config := &mock.MockConfiguration{}
	meta := &mock.MockEventData{Labels: []string{"size/xs"}, ChangedFiles: 8}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "added labels")
	assert.Equal(t, "size/s", meta.AddedLabels[0], "size label")
	assert.Len(t, meta.RemovedLabels, 1, "removed labels")
	assert.Equal(t, "size/xs", meta.RemovedLabels[0], "size label")
}

func TestSizingNoUpdateNeeded(t *testing.T) {
	action := buildRules(false)
	config := &mock.MockConfiguration{}
	meta := &mock.MockEventData{Labels: []string{"size/xs"}, ChangedFiles: 2}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 0, "added labels")
	assert.Len(t, meta.RemovedLabels, 0, "removed labels")
}

func TestSizingNoChanges(t *testing.T) {
	action := buildRules(false)
	config := &mock.MockConfiguration{}
	meta := &mock.MockEventData{ChangedFiles: 0, ChangesAdd: 0, ChangesRemove: 0}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 0, "added labels")
	assert.Len(t, meta.RemovedLabels, 0, "removed labels")
}
