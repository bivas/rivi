package api

import "github.com/bivas/rivi/types"

type JobHandler interface {
	Send(types.Data)
}
