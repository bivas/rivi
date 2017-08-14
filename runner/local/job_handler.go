package local

import (
	"github.com/bivas/rivi/runner/internal"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
)

type jobHandler struct {
	channel chan types.HookData
	work    *workUnit

	logger log.Logger
}

func (h *jobHandler) Send(data types.HookData) {
	h.channel <- data
}

func (h *jobHandler) Start() {
	go h.work.Handle()
}

func NewJobHandler() internal.JobHandler {
	c := make(chan types.HookData)
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
