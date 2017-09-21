package commenter

import (
	"fmt"

	"github.com/bivas/rivi/engine/actions"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/multistep"
)

type action struct {
	rule *rule
}

func (a *action) String() string {
	return fmt.Sprintf("%T{rule: %+v}", *a, a.rule)
}

func (a *action) Apply(state multistep.StateBag) {
	state.Get("data").(types.Data).AddComment(a.rule.Comment)
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) actions.Action {
	item := rule{}
	if e := mapstructure.Decode(config, &item); e != nil {
		log.Get("commenter").ErrorWith(log.MetaFields{log.E(e)}, "Unable to build action")
		return nil
	}
	return &action{rule: &item}
}

func init() {
	actions.RegisterAction("commenter", &factory{})
}
