package bot

import (
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
)

type HookListenerQueue interface {
	Send(types.Data)
}

type channelHookListenerQueue struct {
	incoming chan types.Data
}

func (c *channelHookListenerQueue) Send(data types.Data) {
	c.incoming <- data
}

type HookListenerQueueCreator func() HookListenerQueue

func channelHookListenerQueueCreator() HookListenerQueue {
	log.Get("hook.listener.channel").Debug("Creating hook listener queue creator")
	incomingHooks := make(chan types.Data)
	go runHookHandler(incomingHooks)
	return &channelHookListenerQueue{incomingHooks}
}

var defaultHookListenerQueueCreator = channelHookListenerQueueCreator

func SetHookListenerQueueCreator(fn HookListenerQueueCreator) {
	defaultHookListenerQueueCreator = fn
}

func CreateHookListenerQueue() HookListenerQueue {
	return defaultHookListenerQueueCreator()
}
