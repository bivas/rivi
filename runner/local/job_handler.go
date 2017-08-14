package local

import (
	"github.com/bivas/rivi/runner/api"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
)

type jobHandler struct {
	channel chan types.Data
	work    *workUnit

	logger log.Logger
}

func (h *jobHandler) Send(data types.Data) {
	h.channel <- data
}

func (h *jobHandler) Start() {
	go h.work.Handle()
}

func NewJobHandler() api.JobHandler {
	c := make(chan types.Data)
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
