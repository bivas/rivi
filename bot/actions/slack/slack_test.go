package slack

import (
	"github.com/bivas/rivi/bot/mock"
	"testing"
)

func TestActionApplyRequire1NoAssignees(t *testing.T) {
	action := action{rule: &rule{}}
	config := &mock.MockConfiguration{}
	meta := &mock.MockEventData{}
	action.Apply(config, meta)
	panic("implement me")
}
