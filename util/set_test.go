package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	values := []string{"one", "one", "two"}
	set := StringSet{}
	set.AddAll(values)
	result := set.Values()
	assert.Len(t, result, 2, "2 items")
	assert.Contains(t, result, "one", "one")
	assert.Contains(t, result, "two", "two")
}

func TestSetWithTransform(t *testing.T) {
	values := []string{"ONE", "ONE", "TWO"}
	set := StringSet{Transformer: strings.ToLower}
	set.AddAll(values)
	result := set.Values()
	assert.Len(t, result, 2, "2 items")
	assert.Contains(t, result, "one", "one")
	assert.Contains(t, result, "two", "two")
}
