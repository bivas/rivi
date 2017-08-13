package bot

import (
	"github.com/bivas/rivi/bot/api"
	"github.com/bivas/rivi/bot/local"
)

var defaultHookListenerQueueProvider api.HookListenerQueueProvider = local.CreateHookListenerQueue

func SetHookListenerQueueProvider(fn api.HookListenerQueueProvider) {
	defaultHookListenerQueueProvider = fn
}

func CreateHookListenerQueue() api.HookListenerQueue {
	return defaultHookListenerQueueProvider()
}
