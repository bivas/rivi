package bot

import (
	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/spf13/viper"
	"net/http"
)

type hookListener struct {
	conf  client.ClientConfig
	queue HookListenerQueue

	logger log.Logger
}

func (h *hookListener) HandleEvent(r *http.Request) *HandledEventResult {
	data, ok := types.BuildFromHook(h.conf, r)
	if !ok {
		return &HandledEventResult{Message: "Skipping hook processing"}
	}
	h.queue.Send(data)
	return &HandledEventResult{Message: "Processing " + data.GetShortName()}
}

func NewHookListener() (Bot, error) {
	logger := log.Get("hook.listener")
	return &hookListener{
		conf:   client.NewClientConfig(viper.New()),
		queue:  CreateHookListenerQueue(),
		logger: logger,
	}, nil
}
