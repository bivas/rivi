package local

import (
	"github.com/bivas/rivi/pkg/runner/types"
	"github.com/bivas/rivi/pkg/util/log"
)

type jobHandler struct {
	channel chan *types.Message
	work    *workUnit

	logger log.Logger
}

func (h *jobHandler) Send(data *types.Message) {
	if data == nil {
		h.logger.Debug("Closing channel")
		close(h.channel)
		return
	}
	h.channel <- data
}

func (h *jobHandler) Start() {
	go h.work.Handle()
}

func NewJobHandler() types.JobHandler {
	c := make(chan *types.Message)
	h := &jobHandler{
		channel: c,
		work: &workUnit{
			incoming: c,
			logger:   log.Get("workunit.local")},
		logger: log.Get("job.handler.local"),
	}
	h.Start()
	return h
}
