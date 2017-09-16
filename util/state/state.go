package state

import (
	"github.com/bivas/rivi/config"
	"github.com/bivas/rivi/types"
	"github.com/mitchellh/multistep"
)

func New(config config.Configuration, data types.Data) multistep.StateBag {
	result := new(multistep.BasicStateBag)
	result.Put("data", data)
	result.Put("config", config)
	return result
}
