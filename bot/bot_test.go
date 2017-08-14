package bot

import (
	"net/http"
	"testing"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/mocks"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/types/builder"
	"github.com/stretchr/testify/assert"
	"strings"
)

type mockDataBuilder struct {
	Labels   []string
	Provider string
}

func (m *mockDataBuilder) BuildFromPayload(config client.ClientConfig, payload []byte) (types.Data, bool, error) {
	return &mocks.MockData{Labels: m.Labels, Provider: strings.ToLower(m.Provider)}, true, nil
}

func (m *mockDataBuilder) BuildFromHook(config client.ClientConfig, r *http.Request) (types.Data, bool, error) {
	return &mocks.MockData{Labels: m.Labels, Provider: strings.ToLower(m.Provider)}, true, nil
}

func buildRequest(t *testing.T, url string) *http.Request {
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("Error while building request. %s", err)
	}
	return request
}

func TestNewBotDefaultNamespace(t *testing.T) {
	builder.RegisterNewDataBuilder("TestNewBotDefaultNamespace",
		&mockDataBuilder{Provider: "TestNewBotDefaultNamespace", Labels: []string{}})
	b, err := New("../config/config_test.yml", "../config/empty_config_test.yml")
	if err != nil {
		t.Fatalf("Error while building a bot. %s", err)
	}
	request := buildRequest(t, "http://localhost/")
	request.Header.Add("X-TestNewBotDefaultNamespace", "mock")
	response := b.HandleEvent(request)
	assert.Len(t, response.AppliedRules, 1, "no rules applied")
	assert.Equal(t, "rule2", response.AppliedRules[0], "rule2 on default")
	assert.Empty(t, response.Message, "has error")
}

func TestNewBotExistingNamespace(t *testing.T) {
	builder.RegisterNewDataBuilder("TestNewBotExistingNamespace",
		&mockDataBuilder{
			Provider: "TestNewBotExistingNamespace",
			Labels: []string{
				"label2",
				"pending-approval"},
		})
	b, err := New("../config/config_test.yml", "../config/empty_config_test.yml")
	if err != nil {
		t.Fatalf("Error while building a bot. %s", err)
	}
	request := buildRequest(t, "http://localhost/?namespace=empty_config_test")
	if err != nil {
		t.Fatalf("Error while building request. %s", err)
	}
	request.Header.Add("X-TestNewBotExistingNamespace", "mock")
	response := b.HandleEvent(request)
	assert.Len(t, response.AppliedRules, 1, "no rules applied")
	assert.Equal(t, "rule1", response.AppliedRules[0], "rule1 on namespace")
	assert.Empty(t, response.Message, "has error")
}

func TestNewBotNonExistingNamespace(t *testing.T) {
	builder.RegisterNewDataBuilder("TestNewBotNonExistingNamespace",
		&mockDataBuilder{
			Labels: []string{},
		})
	b, err := New("../config/config_test.yml", "../config/empty_config_test.yml")
	if err != nil {
		t.Fatalf("Error while building a bot. %s", err)
	}
	request := buildRequest(t, "http://localhost/?namespace=TestNewBotNonExistingNamespace")
	if err != nil {
		t.Fatalf("Error while building request. %s", err)
	}
	request.Header.Add("X-TestNewBotNonExistingNamespace", "mock")
	response := b.HandleEvent(request)
	assert.Len(t, response.AppliedRules, 0, "no rules applied")
	assert.NotEmpty(t, response.Message, "should have error")
}
