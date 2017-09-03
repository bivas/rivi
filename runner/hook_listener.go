package runner

import (
	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/runner/internal"
	"github.com/bivas/rivi/types/builder"
	"github.com/bivas/rivi/util/log"
	"github.com/spf13/viper"
	"net/http"
)

type hookListener struct {
	conf  client.ClientConfig
	queue internal.HookListenerQueue

	logger log.Logger
}

func (h *hookListener) HandleEvent(r *http.Request) *HandledEventResult {
	data, ok := builder.BuildFromHook(h.conf, r)
	if !ok {
		return &HandledEventResult{Message: "Skipping hook processing"}
	}
	h.queue.Enqueue(data)
	return &HandledEventResult{Message: "Processing " + data.GetShortName()}
}

func NewHookListener() (Runner, error) {
	logger := runnerLog.Get("hook.listener")
	return &hookListener{
		conf:   client.NewClientConfig(viper.New()),
		queue:  CreateHookListenerQueue(),
		logger: logger,
	}, nil
}
