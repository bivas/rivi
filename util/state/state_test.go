package state

import (
	"github.com/bivas/rivi/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestState(t *testing.T) {
	meta := &mocks.MockData{}
	config := &mocks.MockConfiguration{}

	result := New(config, meta)
	assert.Equal(t, meta, result.Get("data"), "data")
	assert.Equal(t, config, result.Get("config"), "config")
}
