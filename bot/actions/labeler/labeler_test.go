package labeler

import (
	"testing"

	"github.com/bivas/rivi/bot/mock"
	"github.com/bivas/rivi/util/log"
	"github.com/stretchr/testify/assert"
)

func TestSerialization(t *testing.T) {
	input := map[string]interface{}{
		"label":  "label1",
		"remove": "label2",
	}

	var f factory
	result := f.BuildAction(input)
	assert.NotNil(t, result, "should create action")
	s, ok := result.(*action)
	assert.True(t, ok, "should be of this package")
	assert.Equal(t, "label1", s.rule.Label, "label")
	assert.Equal(t, "label2", s.rule.Remove, "remove")
}

func TestLabelerNoLabels(t *testing.T) {
	action := action{rule: &rule{Label: "label1"}, logger: log.Get("labeler.test")}
	meta := &mock.MockEventData{Labels: []string{}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "added labels")
	assert.Equal(t, meta.Labels, []string{"label1"}, "labels")
}

func TestLabelExists(t *testing.T) {
	action := action{rule: &rule{Label: "label1"}, logger: log.Get("labeler.test")}
	meta := &mock.MockEventData{Labels: []string{"label1"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 0, "added labels")
}

func TestNewLabelWithExisting(t *testing.T) {
	action := action{rule: &rule{Label: "label1"}, logger: log.Get("labeler.test")}
	meta := &mock.MockEventData{Labels: []string{"label2"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "added labels")
	assert.Len(t, meta.Labels, 2, "labels")
}

func TestRemoveNotExisting(t *testing.T) {
	action := action{rule: &rule{Label: "label2", Remove: "label1"}, logger: log.Get("labeler.test")}
	meta := &mock.MockEventData{Labels: []string{"label2"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.RemovedLabels, 0, "removed labels")
	assert.Len(t, meta.Labels, 1, "labels")
}

func TestRemoveExisting(t *testing.T) {
	action := action{rule: &rule{Remove: "label1"}, logger: log.Get("labeler.test")}
	meta := &mock.MockEventData{Labels: []string{"label1"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.RemovedLabels, 1, "removed labels")
	assert.Len(t, meta.Labels, 0, "labels")
}

func TestReplaceLabeles(t *testing.T) {
	action := action{rule: &rule{Remove: "label1", Label: "label2"}, logger: log.Get("labeler.test")}
	meta := &mock.MockEventData{Labels: []string{"label1"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.RemovedLabels, 1, "removed labels")
	assert.Len(t, meta.AddedLabels, 1, "added labels")
	assert.Len(t, meta.Labels, 1, "labels")
	assert.Equal(t, meta.Labels[0], "label2", "label2")
}
