package runner

import (
	"github.com/bivas/rivi/runner/internal"
	"github.com/bivas/rivi/runner/local"
)

var defaultHookListenerQueueProvider internal.HookListenerQueueProvider = local.CreateHookListenerQueue

func SetHookListenerQueueProvider(fn internal.HookListenerQueueProvider) {
	defaultHookListenerQueueProvider = fn
}

func CreateHookListenerQueue() internal.HookListenerQueue {
	return defaultHookListenerQueueProvider()
}
