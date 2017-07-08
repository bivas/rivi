package locker

import (
	"github.com/bivas/rivi/bot/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLabelerNoLabels(t *testing.T) {
	action := action{rule: &rule{}}
	meta := &mock.MockEventData{Labels: []string{}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.True(t, true)
}
