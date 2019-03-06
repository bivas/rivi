package local

import (
	"github.com/bivas/rivi/pkg/runner/types"
	"github.com/bivas/rivi/pkg/util/log"
)

type channelHookListenerQueue struct {
	incoming chan *types.Message
}

func (c *channelHookListenerQueue) Enqueue(data *types.Message) {
	if data == nil {
		close(c.incoming)
		return
	}
	c.incoming <- data
}

func channelHookListenerQueueProvider() types.HookListenerQueue {
	log.Get("hook.listener.queue").DebugWith(
		log.MetaFields{
			log.F("type", "channel")}, "Creating hook listener queue provider")
	incomingHooks := make(chan *types.Message)
	handler := NewChannelHookHandler(incomingHooks)
	go handler.Run()
	return &channelHookListenerQueue{incomingHooks}
}

func CreateHookListenerQueue() types.HookListenerQueue {
	return channelHookListenerQueueProvider()
}
