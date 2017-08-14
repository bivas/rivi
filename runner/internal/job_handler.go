package internal

import "github.com/bivas/rivi/types"

type JobHandler interface {
	Send(types.HookData)
}
