package local

import (
	"time"

	"github.com/bivas/rivi/bot/api"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"

	"github.com/patrickmn/go-cache"
)

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
				"Creating new job handler")
			c = NewJobHandler()
		}
		h.processingCache.Set(key, c, cache.DefaultExpiration)
		h.logger.DebugWith(
			log.MetaFields{
				log.F("issue", key),
			}, "Sending data to job handler")
		c.(api.JobHandler).Send(data)
	}
}

func NewChannelHookHandler(incoming <-chan types.Data) api.HookHandler {
	return &channelHookHandler{
		incoming:        incoming,
		processingCache: cache.New(time.Minute, 2*time.Minute),
		logger:          log.Get("hook.handler")}
}
