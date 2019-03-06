package runner

import (
	"net/http"

	"github.com/bivas/rivi/pkg/config/client"
	"github.com/bivas/rivi/pkg/runner/types"
	"github.com/bivas/rivi/pkg/types/builder"
	"github.com/bivas/rivi/pkg/util/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

type hookListener struct {
	conf  client.ClientConfig
	queue types.HookListenerQueue

	logger log.Logger
}

func (h *hookListener) HandleEvent(r *http.Request) *HandledEventResult {
	timer := prometheus.NewTimer(incomingWebhooksHistogram)
	defer timer.ObserveDuration()
	data, ok := builder.BuildFromHook(h.conf, r)
	if !ok {
		skippedWebhooksCounter.Inc()
		return &HandledEventResult{Message: "Skipping hook processing"}
	}
	h.queue.Enqueue(types.NewMessage(h.conf, data))
	return &HandledEventResult{Message: "Processing " + data.GetShortName()}
}

func NewHookListener(clientConfiguration string) (Runner, error) {
	var conf client.ClientConfig
	if clientConfiguration == "" {
		conf = client.NewClientConfig(viper.New())
	} else {
		conf = client.NewClientConfigFromFile(clientConfiguration)
	}
	return NewHookListenerWithClientConfig(conf), nil
}

func NewHookListenerWithClientConfig(config client.ClientConfig) Runner {
	logger := runnerLog.Get("hook.listener")
	return &hookListener{
		conf:   config,
		queue:  CreateHookListenerQueue(),
		logger: logger,
	}
}

var incomingWebhooksHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
	Namespace: "rivi",
	Subsystem: "webhook",
	Name:      "incoming",
	Help:      "Measure incoming webhook processing",
})

var skippedWebhooksCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "rivi",
	Subsystem: "webhook",
	Name:      "skipped",
	Help:      "Measure skipped webhooks",
})

func init() {
	prometheus.Register(incomingWebhooksHistogram)
	prometheus.Register(skippedWebhooksCounter)
}
