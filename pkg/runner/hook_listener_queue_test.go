package runner

import (
	"testing"

	"github.com/bivas/rivi/pkg/runner/types"
	"github.com/stretchr/testify/assert"
)

type testHookListenerQueue struct {
	Created bool
}

func (*testHookListenerQueue) Enqueue(*types.Message) {
	panic("implement me")
}

func TestSetHookListenerQueueProvider(t *testing.T) {
	queue := &testHookListenerQueue{false}
	SetHookListenerQueueProvider(func() types.HookListenerQueue {
		queue.Created = true
		return queue
	})
	assert.False(t, queue.Created, "created should be false")
	result := CreateHookListenerQueue()
	assert.NotNil(t, result, "test result")
	_, ok := result.(*testHookListenerQueue)
	assert.True(t, ok, "same type")
	assert.True(t, queue.Created, "created")
}
