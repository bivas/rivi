package runner

import (
	"net/http"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/runner/internal"
	"github.com/bivas/rivi/types/builder"
	"github.com/bivas/rivi/util/log"
	"github.com/spf13/viper"
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
	h.queue.Enqueue(internal.NewMessage(h.conf, data))
	return &HandledEventResult{Message: "Processing " + data.GetShortName()}
}

func NewHookListener(clientConfiguration string) (Runner, error) {
	logger := runnerLog.Get("hook.listener")
	var conf client.ClientConfig
	if clientConfiguration == "" {
		conf = client.NewClientConfig(viper.New())
	} else {
		conf = client.NewClientConfigFromFile(clientConfiguration)
	}
	return &hookListener{
		conf:   conf,
		queue:  CreateHookListenerQueue(),
		logger: logger,
	}, nil
}
