package internal

import (
	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/types"
)

type Message struct {
	Config client.ClientConfig
	Data   types.HookData
}
