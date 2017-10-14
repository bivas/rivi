package runner

import (
	"net/http"
	"testing"

	"strings"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/mocks"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/types/builder"
	"github.com/stretchr/testify/assert"
)

type mockDataBuilder struct {
	Labels   []string
	Provider string
}

func (m *mockDataBuilder) BuildFromPayload(config client.ClientConfig, payload []byte) (types.Data, bool, error) {
	return &mocks.MockData{Labels: m.Labels, Provider: strings.ToLower(m.Provider)}, true, nil
}

func (m *mockDataBuilder) BuildFromHook(config client.ClientConfig, r *http.Request) (types.HookData, bool, error) {
	return &mocks.MockData{Labels: m.Labels, Provider: strings.ToLower(m.Provider)}, true, nil
}

func buildRequest(t *testing.T, url string) *http.Request {
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("Error while building request. %s", err)
	}
	return request
}

func TestNewRunnerDefaultNamespace(t *testing.T) {
	builder.RegisterNewDataBuilder("TestNewRunnerDefaultNamespace",
		&mockDataBuilder{Provider: "TestNewRunnerDefaultNamespace", Labels: []string{}})
	b, err := New("../config/config_test.yml", "../config/empty_config_test.yml")
	if err != nil {
		t.Fatalf("Error while building a runnable. %s", err)
	}
	request := buildRequest(t, "http://localhost/")
	request.Header.Set("User-Agent", "X-TestNewRunnerDefaultNamespace")
	response := b.HandleEvent(request)
	assert.Len(t, response.AppliedRules, 1, "no rules applied")
	assert.Equal(t, "rule2", response.AppliedRules[0], "rule2 on default")
	assert.Empty(t, response.Message, "has error")
}

func TestNewRunnerExistingNamespace(t *testing.T) {
	builder.RegisterNewDataBuilder("TestNewRunnerExistingNamespace",
		&mockDataBuilder{
			Provider: "TestNewRunnerExistingNamespace",
			Labels: []string{
				"label2",
				"pending-approval"},
		})
	b, err := New("../config/config_test.yml", "../config/empty_config_test.yml")
	if err != nil {
		t.Fatalf("Error while building a runnable. %s", err)
	}
	request := buildRequest(t, "http://localhost/?namespace=empty_config_test")
	if err != nil {
		t.Fatalf("Error while building request. %s", err)
	}
	request.Header.Set("User-Agent", "X-TestNewRunnerExistingNamespace")
	response := b.HandleEvent(request)
	assert.Len(t, response.AppliedRules, 1, "no rules applied")
	assert.Equal(t, "rule1", response.AppliedRules[0], "rule1 on namespace")
	assert.Empty(t, response.Message, "has error")
}

func TestNewRunnerNonExistingNamespace(t *testing.T) {
	builder.RegisterNewDataBuilder("TestNewRunnerNonExistingNamespace",
		&mockDataBuilder{
			Labels: []string{},
		})
	b, err := New("../config/config_test.yml", "../config/empty_config_test.yml")
	if err != nil {
		t.Fatalf("Error while building a runnable. %s", err)
	}
	request := buildRequest(t, "http://localhost/?namespace=TestNewRunnerNonExistingNamespace")
	if err != nil {
		t.Fatalf("Error while building request. %s", err)
	}
	request.Header.Set("User-Agent", "X-TestNewRunnerNonExistingNamespace")
	response := b.HandleEvent(request)
	assert.Len(t, response.AppliedRules, 0, "no rules applied")
	assert.NotEmpty(t, response.Message, "should have error")
}
