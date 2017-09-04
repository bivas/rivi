package local

import (
	"time"

	"github.com/bivas/rivi/runner/internal"
	"github.com/bivas/rivi/util/log"

	"github.com/patrickmn/go-cache"
)

type channelHookHandler struct {
	incoming <-chan internal.Message
	logger   log.Logger

	processingCache *cache.Cache
}

func (h *channelHookHandler) Run() {
	for {
		msg, ok := <-h.incoming
		if !ok {
			h.logger.Info("Hook channel has no more data - exiting")
			break
		}
		key := msg.Data.GetShortName()
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
		c.(internal.JobHandler).Send(msg)
	}
}

func NewChannelHookHandler(incoming <-chan internal.Message) internal.HookHandler {
	return &channelHookHandler{
		incoming:        incoming,
		processingCache: cache.New(time.Minute, 2*time.Minute),
		logger:          log.Get("hook.handler")}
}
