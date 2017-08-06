package bot

import (
	"net/http"
	"testing"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/mocks"
	"github.com/bivas/rivi/types"
	"github.com/stretchr/testify/assert"
)

type mockEventDataBuilder struct {
	Labels []string
}

func (m *mockEventDataBuilder) BuildFromPayload(config client.ClientConfig, payload []byte) (types.EventData, bool, error) {
	return &mocks.MockEventData{Labels: m.Labels}, true, nil
}

func (m *mockEventDataBuilder) BuildFromRequest(config client.ClientConfig, r *http.Request) (types.EventData, bool, error) {
	return &mocks.MockEventData{Labels: m.Labels}, true, nil
}

func (m *mockEventDataBuilder) PartialBuildFromRequest(config client.ClientConfig, r *http.Request) (types.EventData, bool, error) {
	return &mocks.MockEventData{Labels: m.Labels}, true, nil
}

func buildRequest(t *testing.T, url string) *http.Request {
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("Error while building request. %s", err)
	}
	return request
}

func TestNewBotDefaultNamespace(t *testing.T) {
	types.RegisterNewDataBuilder("TestNewBotDefaultNamespace",
		&mockEventDataBuilder{Labels: []string{}})
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
	types.RegisterNewDataBuilder("TestNewBotExistingNamespace",
		&mockEventDataBuilder{
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
	types.RegisterNewDataBuilder("TestNewBotNonExistingNamespace",
		&mockEventDataBuilder{
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
