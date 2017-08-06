package bot

import (
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
)

type JobHandler interface {
	Handle(chan types.Data)
}

type loggerJobHandler struct {
	logger log.Logger
}

func (h *loggerJobHandler) Handle(incoming chan types.Data) {
	for {
		data, ok := <-incoming
		if !ok {
			h.logger.Info("Stopping job handler")
			break
		}
		log.InfoWith(log.MetaFields{log.F("data", data.GetShortName())}, "Got data from job channel")
	}
}
