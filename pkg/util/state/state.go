package state

import (
	"github.com/bivas/rivi/pkg/config"
	"github.com/bivas/rivi/pkg/types"
	"github.com/mitchellh/multistep"
)

func New(config config.Configuration, data types.Data) multistep.StateBag {
	result := new(multistep.BasicStateBag)
	result.Put("data", data)
	result.Put("config", config)
	return result
}
