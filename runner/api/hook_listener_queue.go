package api

import (
	"github.com/bivas/rivi/types"
)

type HookListenerQueue interface {
	Send(types.Data)
}

type HookListenerQueueProvider func() HookListenerQueue
