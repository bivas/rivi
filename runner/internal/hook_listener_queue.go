package internal

import (
	"github.com/bivas/rivi/types"
)

type HookListenerQueue interface {
	Send(types.HookData)
}

type HookListenerQueueProvider func() HookListenerQueue
