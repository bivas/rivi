package sizing

import (
	"github.com/bivas/rivi/bot/mock"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func buildRule(name, label string, changedFiles int) *sizingRule {
	result := &sizingRule{Name: name, Label: label, ChangedFilesThreshold: changedFiles}
	result.Defaults()
	return result
}

func buildRules(withDefault bool) *action {
	result := &action{items: rules{
		*buildRule("xs", "size/xs", 5),
		*buildRule("s", "size/s", 15),
		*buildRule("m", "size/m", 75),
	}, possibleLabels: []string{
		"size/xs",
		"size/s",
		"size/m",
	}}
	if withDefault {
		result.possibleLabels = append(result.possibleLabels, "default-label")
		result.items = append(result.items, sizingRule{Name: "default", Label: "default-label"})
	}
	sort.Sort(sort.Reverse(result.items))
	return result
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
