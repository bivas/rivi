package internal

import (
	"github.com/bivas/rivi/types"
)

type HookListenerQueue interface {
	Enqueue(types.HookData)
}

type HookListenerQueueProvider func() HookListenerQueue
