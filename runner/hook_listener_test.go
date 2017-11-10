package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHookListener(t *testing.T) {
	result, err := NewHookListener("")
	assert.NoError(t, err, "new failed")
	assert.NotNil(t, result, "nil runner")
	_, ok := result.(*hookListener)
	assert.True(t, ok, "must be of type hook lisetner")
}
