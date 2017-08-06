package bot

import (
	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/spf13/viper"
	"net/http"
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

type hookListener struct {
	conf  client.ClientConfig
	queue HookListenerQueue

	logger log.Logger
}

func (h *hookListener) HandleEvent(r *http.Request) *HandledEventResult {
	data, ok := types.BuildFromHook(h.conf, r)
	if !ok {
		return &HandledEventResult{Message: "Skipping rules processing"}
	}
	h.queue.Send(data)
	return &HandledEventResult{Message: "Processing " + data.GetLongName()}
}

func NewHookListener() (Bot, error) {
	logger := log.Get("hooklistener")
	incomingHooks := make(chan types.Data)
	go runHookHandler(incomingHooks)
	return &hookListener{
		conf:   client.NewClientConfig(viper.New()),
		queue:  &channelHookListenerQueue{incomingHooks},
		logger: logger,
	}, nil
}
