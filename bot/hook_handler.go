package bot

import (
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/patrickmn/go-cache"
	"time"
)

type HookHandler interface {
	Run()
}

type processingUnit struct {
	Channel chan types.Data
	Handler JobHandler

	logger log.Logger
}

func (p *processingUnit) Start() {
	go p.Handler.Handle(p.Channel)
}

type channelHookHandler struct {
	incoming <-chan types.Data
	logger   log.Logger

	processingCache *cache.Cache
}

func (h *channelHookHandler) Run() {
	for {
		data, ok := <-h.incoming
		if !ok {
			h.logger.Info("Hook channel has no more data - exiting")
			break
		}
		key := data.GetShortName()
		c, exists := h.processingCache.Get(key)
		if !exists {
			h.logger.DebugWith(
				log.MetaFields{
					log.F("issue", key),
				},
				"Creating new processing unit")
			c = &processingUnit{
				Channel: make(chan types.Data),
				// abstract this handler creation
				Handler: &localJobHandler{logger: h.logger.Get("job.handler")},
				logger:  h.logger.Get("unit"),
			}
			c.(*processingUnit).Start()
		}
		h.processingCache.Set(key, c, cache.DefaultExpiration)
		h.logger.DebugWith(
			log.MetaFields{
				log.F("issue", key),
				log.F("pending", len(c.(*processingUnit).Channel)),
			}, "Sending data to processing unit")
		c.(*processingUnit).Channel <- data
	}
}

func runHookHandler(incoming <-chan types.Data) {
	handler := &channelHookHandler{
		incoming:        incoming,
		processingCache: cache.New(time.Minute, 2*time.Minute),
		logger:          log.Get("hook.handler")}
	handler.Run()
}
