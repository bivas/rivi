package labeler

import (
	"github.com/bivas/rivi/bot/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLabelerNoLabels(t *testing.T) {
	action := action{rule: &rule{Label:"label1"}}
	meta := &mock.MockEventData{Labels:[]string{}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "added labels")
	assert.Equal(t, meta.Labels, []string{"label1"}, "labels")
}

func TestLabelExists(t *testing.T) {
	action := action{rule: &rule{Label:"label1"}}
	meta := &mock.MockEventData{Labels:[]string{"label1"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 0, "added labels")
}

func TestNewLabelWithExisting(t *testing.T) {
	action := action{rule: &rule{Label:"label1"}}
	meta := &mock.MockEventData{Labels:[]string{"label2"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedLabels, 1, "added labels")
	assert.Len(t, meta.Labels, 2, "labels")
}