package internal

type HookListenerQueue interface {
	Enqueue(Message)
}

type HookListenerQueueProvider func() HookListenerQueue
