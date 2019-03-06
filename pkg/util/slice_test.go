package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInStringSlice(t *testing.T) {
	slice := []string{"a", "b", "c"}
	assert.False(t, InStringSlice(slice, "x"), "x")
	assert.True(t, InStringSlice(slice, "c"), "c")
}
