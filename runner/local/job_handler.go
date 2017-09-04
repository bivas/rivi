package local

import (
	"github.com/bivas/rivi/runner/internal"
	"github.com/bivas/rivi/util/log"
)

type jobHandler struct {
	channel chan internal.Message
	work    *workUnit

	logger log.Logger
}

func (h *jobHandler) Send(data internal.Message) {
	h.channel <- data
}

func (h *jobHandler) Start() {
	go h.work.Handle()
}

func NewJobHandler() internal.JobHandler {
	c := make(chan internal.Message)
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
