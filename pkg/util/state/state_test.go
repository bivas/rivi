package state

import (
	"testing"

	"github.com/bivas/rivi/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	meta := &mocks.MockData{}
	config := &mocks.MockConfiguration{}

	result := New(config, meta)
	assert.Equal(t, meta, result.Get("data"), "data")
	assert.Equal(t, config, result.Get("config"), "config")
}
