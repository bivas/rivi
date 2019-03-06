package types

type HookListenerQueue interface {
	Enqueue(*Message)
}

type HookListenerQueueProvider func() HookListenerQueue
