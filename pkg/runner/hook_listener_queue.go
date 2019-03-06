package runner

import (
	"github.com/bivas/rivi/pkg/runner/local"
	"github.com/bivas/rivi/pkg/runner/types"
)

var defaultHookListenerQueueProvider types.HookListenerQueueProvider = local.CreateHookListenerQueue

func SetHookListenerQueueProvider(fn types.HookListenerQueueProvider) {
	defaultHookListenerQueueProvider = fn
}

func CreateHookListenerQueue() types.HookListenerQueue {
	return defaultHookListenerQueueProvider()
}
