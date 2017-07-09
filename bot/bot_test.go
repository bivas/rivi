package bot

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type mockEventDataBuilder struct {
	Labels []string
}

func (m *mockEventDataBuilder) BuildFromRequest(config ClientConfig, r *http.Request) (EventData, bool, error) {
	return &mockConditionEventData{Labels: m.Labels}, true, nil
}

func (m *mockEventDataBuilder) PartialBuildFromRequest(config ClientConfig, r *http.Request) (EventData, bool, error) {
	return &mockConditionEventData{Labels: m.Labels}, true, nil
}

func buildRequest(t *testing.T, url string) *http.Request {
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("Error while building request. %s", err)
	}
	return request
}

func TestNewBotDefaultNamespace(t *testing.T) {
	RegisterNewBuilder("TestNewBotDefaultNamespace",
		&mockEventDataBuilder{Labels: []string{}})
	b, err := New("config_test.yml", "empty_config_test.yml")
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
	RegisterNewBuilder("TestNewBotExistingNamespace",
		&mockEventDataBuilder{
			Labels: []string{
				"label2",
				"pending-approval"},
		})
	b, err := New("config_test.yml", "empty_config_test.yml")
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
	RegisterNewBuilder("TestNewBotNonExistingNamespace",
		&mockEventDataBuilder{
			Labels: []string{},
		})
	b, err := New("config_test.yml", "empty_config_test.yml")
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
