package github

import (
	"net/http"
	"testing"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/mocks"
	"github.com/bivas/rivi/types"
	"github.com/stretchr/testify/assert"
)

type mockEventHandler struct {
	FromRequestCalled bool
	FromPayloadCalled bool
}

func (m *mockEventHandler) FromRequest(client.ClientConfig, *http.Request) (types.HookData, bool, error) {
	m.FromRequestCalled = true
	return nil, false, nil
}

func (m *mockEventHandler) FromPayload(client.ClientConfig, []byte) (types.Data, bool, error) {
	m.FromPayloadCalled = true
	return nil, false, nil
}

func TestRequestDefault(t *testing.T) {
	DataBuilder.handlers["mock"] = &mockEventHandler{}
	DataBuilder.defaultHandler = &mockEventHandler{}

	request, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err, "shouldn't error")
	_, cont, _ := DataBuilder.BuildFromHook(&mocks.MockClientConfig{}, request)
	assert.False(t, cont, "shouldn't continue")
	assert.True(t, DataBuilder.defaultHandler.(*mockEventHandler).FromRequestCalled)
}

func TestRequest(t *testing.T) {
	DataBuilder.handlers["mock"] = &mockEventHandler{}
	DataBuilder.defaultHandler = &mockEventHandler{}

	request, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err, "shouldn't error")
	request.Header.Set("X-Github-Event", "mock")
	_, cont, _ := DataBuilder.BuildFromHook(&mocks.MockClientConfig{}, request)
	assert.False(t, cont, "shouldn't continue")
	assert.False(t, DataBuilder.defaultHandler.(*mockEventHandler).FromRequestCalled)
	assert.True(t, DataBuilder.handlers["mock"].(*mockEventHandler).FromRequestCalled)
}
