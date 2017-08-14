package local

import (
	"github.com/bivas/rivi/runner/api"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
)

type channelHookListenerQueue struct {
	incoming chan types.Data
}

func (c *channelHookListenerQueue) Send(data types.Data) {
	c.incoming <- data
}

func channelHookListenerQueueProvider() api.HookListenerQueue {
	log.Get("hook.listener.queue").DebugWith(
		log.MetaFields{
			log.F("type", "channel")}, "Creating hook listener queue provider")
	incomingHooks := make(chan types.Data)
	handler := NewChannelHookHandler(incomingHooks)
	go handler.Run()
	return &channelHookListenerQueue{incomingHooks}
}

func CreateHookListenerQueue() api.HookListenerQueue {
	return channelHookListenerQueueProvider()
}
