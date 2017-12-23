package status

import (
	"testing"

	"github.com/bivas/rivi/mocks"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/state"
	"github.com/stretchr/testify/assert"
)

func TestSerialization(t *testing.T) {
	input := map[string]interface{}{
		"state":       "Success",
		"description": "this is a test",
	}

	var f factory
	result := f.BuildAction(input)
	assert.NotNil(t, result, "should create action")
	s, ok := result.(*action)
	assert.True(t, ok, "should be of this package")
	assert.Equal(t, "Success", s.rule.State, "state")
	assert.Equal(t, "this is a test", s.rule.Description, "desc")
}

func TestSetDefaultStatus(t *testing.T) {
	input := map[string]interface{}{
		"description": "TestSetDefaultStatus",
	}
	var f factory
	action := f.BuildAction(input)
	assert.NotNil(t, action, "should create action")
	meta := &mocks.MockData{}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.Equal(t, types.GetState(defaultState), meta.StatusState, "state")
	assert.Equal(t, "TestSetDefaultStatus", meta.StatusDescription, "desc")
}

func TestSetUnknownStatus(t *testing.T) {
	input := map[string]interface{}{
		"state":       "TestSetUnknownStatus",
		"description": "TestSetUnknownStatus",
	}
	var f factory
	action := f.BuildAction(input)
	assert.NotNil(t, action, "should create action")
	meta := &mocks.MockData{}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.Equal(t, types.GetState(unknownState), meta.StatusState, "state")
	assert.Equal(t, "TestSetUnknownStatus", meta.StatusDescription, "desc")
}

func TestSetSuccessStatus(t *testing.T) {
	input := map[string]interface{}{
		"state":       "Success",
		"description": "TestSetSuccessStatus",
	}
	var f factory
	action := f.BuildAction(input)
	assert.NotNil(t, action, "should create action")
	meta := &mocks.MockData{}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.Equal(t, types.Success, meta.StatusState, "state")
	assert.Equal(t, "TestSetSuccessStatus", meta.StatusDescription, "desc")
}

func TestDefault(t *testing.T) {
	input := map[string]interface{}{}
	var f factory
	action := f.BuildAction(input)
	assert.NotNil(t, action, "should create action")
	meta := &mocks.MockData{}
	config := &mocks.MockConfiguration{}
	action.Apply(state.New(config, meta))
	assert.Equal(t, types.GetState(defaultState), meta.StatusState, "state")
	assert.Equal(t, defaultDescription, meta.StatusDescription, "desc")
}
