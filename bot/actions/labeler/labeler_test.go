package labeler

import (
	"testing"

	"github.com/bivas/rivi/bot/mock"
	"github.com/stretchr/testify/assert"
)

func TestLabelerNoLabels(t *testing.T) {
	action := action{rule: &rule{Label: "label1"}}
	meta := &mock.MockEventData{Labels: []string{}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "added labels")
	assert.Equal(t, meta.Labels, []string{"label1"}, "labels")
}

func TestLabelExists(t *testing.T) {
	action := action{rule: &rule{Label: "label1"}}
	meta := &mock.MockEventData{Labels: []string{"label1"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 0, "added labels")
}

func TestNewLabelWithExisting(t *testing.T) {
	action := action{rule: &rule{Label: "label1"}}
	meta := &mock.MockEventData{Labels: []string{"label2"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "added labels")
	assert.Len(t, meta.Labels, 2, "labels")
}

func TestRemoveNotExisting(t *testing.T) {
	action := action{rule: &rule{Label: "label2", Remove: "label1"}}
	meta := &mock.MockEventData{Labels: []string{"label2"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.RemovedLabels, 0, "removed labels")
	assert.Len(t, meta.Labels, 1, "labels")
}

func TestRemoveExisting(t *testing.T) {
	action := action{rule: &rule{Remove: "label1"}}
	meta := &mock.MockEventData{Labels: []string{"label1"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.RemovedLabels, 1, "removed labels")
	assert.Len(t, meta.Labels, 0, "labels")
}

func TestReplaceLabeles(t *testing.T) {
	action := action{rule: &rule{Remove: "label1", Label: "label2"}}
	meta := &mock.MockEventData{Labels: []string{"label1"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.RemovedLabels, 1, "removed labels")
	assert.Len(t, meta.AddedLabels, 1, "added labels")
	assert.Len(t, meta.Labels, 1, "labels")
	assert.Equal(t, meta.Labels[0], "label2", "label2")
}
