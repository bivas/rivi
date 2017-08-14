package runner

import (
	"github.com/bivas/rivi/runner/api"
	"github.com/bivas/rivi/runner/local"
)

var defaultHookListenerQueueProvider api.HookListenerQueueProvider = local.CreateHookListenerQueue

func SetHookListenerQueueProvider(fn api.HookListenerQueueProvider) {
	defaultHookListenerQueueProvider = fn
}

func CreateHookListenerQueue() api.HookListenerQueue {
	return defaultHookListenerQueueProvider()
}
