package local

import (
	"github.com/bivas/rivi/runner/internal"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
)

type channelHookListenerQueue struct {
	incoming chan types.Data
}

func (c *channelHookListenerQueue) Send(data types.Data) {
	c.incoming <- data
}

func channelHookListenerQueueProvider() internal.HookListenerQueue {
	log.Get("hook.listener.queue").DebugWith(
		log.MetaFields{
			log.F("type", "channel")}, "Creating hook listener queue provider")
	incomingHooks := make(chan types.Data)
	handler := NewChannelHookHandler(incomingHooks)
	go handler.Run()
	return &channelHookListenerQueue{incomingHooks}
}

func CreateHookListenerQueue() internal.HookListenerQueue {
	return channelHookListenerQueueProvider()
}
