package local

import (
	"github.com/bivas/rivi/runner/internal"
	"github.com/bivas/rivi/util/log"
)

type channelHookListenerQueue struct {
	incoming chan *internal.Message
}

func (c *channelHookListenerQueue) Enqueue(data *internal.Message) {
	if data == nil {
		close(c.incoming)
		return
	}
	c.incoming <- data
}

func channelHookListenerQueueProvider() internal.HookListenerQueue {
	log.Get("hook.listener.queue").DebugWith(
		log.MetaFields{
			log.F("type", "channel")}, "Creating hook listener queue provider")
	incomingHooks := make(chan *internal.Message)
	handler := NewChannelHookHandler(incomingHooks)
	go handler.Run()
	return &channelHookListenerQueue{incomingHooks}
}

func CreateHookListenerQueue() internal.HookListenerQueue {
	return channelHookListenerQueueProvider()
}
