package action

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	name = "test-action-config"
)

type testActionConfig struct {
}

func (*testActionConfig) Name() string {
	return name
}

type testActionConfigBuilder struct {
	Result *testActionConfig
	Error  error
}

func (t *testActionConfigBuilder) Build(config map[string]interface{}) (ActionConfig, error) {
	return t.Result, t.Error
}

func reset() {
	delete(actionConfigBuilders, name)
}

func TestBuildActionConfig(t *testing.T) {
	reset()
	RegisterActionConfigBuilder(name, &testActionConfigBuilder{&testActionConfig{}, nil})
	v := viper.New()
	v.Set(name, map[string]interface{}{"one": []int{1}, "two": []int{2, 3}})

	result := BuildActionConfigs(v)
	assert.Contains(t, result, name, "map")
}

func TestBuildActionConfigWithError(t *testing.T) {
	reset()
	RegisterActionConfigBuilder(name, &testActionConfigBuilder{&testActionConfig{}, errors.New("test")})
	v := viper.New()
	v.Set(name, map[string]interface{}{"one": []int{1}, "two": []int{2, 3}})

	result := BuildActionConfigs(v)
	assert.Len(t, result, 0, "should be empty")
}

func TestBuildActionConfigWithNoSection(t *testing.T) {
	reset()
	RegisterActionConfigBuilder(name, &testActionConfigBuilder{&testActionConfig{}, nil})
	v := viper.New()

	result := BuildActionConfigs(v)
	assert.Len(t, result, 0, "should be empty")
}
