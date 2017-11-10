package actions

import (
	"github.com/mitchellh/multistep"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	name = "test-action"
)

type testAction struct {
}

func (*testAction) Apply(multistep.StateBag) {
}

type testActionFactory struct {
	Result Action
}

func (t *testActionFactory) BuildAction(config map[string]interface{}) Action {
	return t.Result
}

func reset() {
	delete(registry, name)
}

func TestBuildActionsFromConfiguration(t *testing.T) {
	reset()
	RegisterAction(name, &testActionFactory{&testAction{}})
	v := viper.New()
	v.Set(name, map[string]interface{}{"test-key": 1})

	result := BuildActionsFromConfiguration(v)
	assert.Len(t, result, 1, "result should have an action")
}

func TestBuildActionsFromConfigurationWithCondition(t *testing.T) {
	reset()
	RegisterAction(name, &testActionFactory{&testAction{}})
	v := viper.New()
	v.Set("condition", map[string]interface{}{"test-key": 1})

	result := BuildActionsFromConfiguration(v)
	assert.Len(t, result, 0, "result should not have an action")
}
