package state

import (
	"github.com/bivas/rivi/config"
	"github.com/bivas/rivi/types"
	"github.com/mitchellh/multistep"
)

func New(config config.Configuration, data types.EventData) multistep.StateBag {
	state := new(multistep.BasicStateBag)
	state.Put("data", data)
	state.Put("config", config)
	return state
}
