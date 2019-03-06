package types

import (
	"github.com/bivas/rivi/pkg/config/client"
	"github.com/bivas/rivi/pkg/types"
)

type Message struct {
	Config client.ClientConfig
	Data   types.HookData
}

func NewMessage(config client.ClientConfig, data types.HookData) *Message {
	return &Message{config, data}
}
