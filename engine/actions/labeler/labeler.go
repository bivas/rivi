package labeler

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

	logger log.Logger
}

func (a *action) String() string {
	return fmt.Sprintf("%T{rule: %+v}", *a, a.rule)
}

func (a *action) Apply(state multistep.StateBag) {
	meta := state.Get("data").(types.Data)
	apply := a.rule.Label
	if apply != "" {
		if meta.HasLabel(apply) {
			a.logger.DebugWith(
				log.MetaFields{log.F("issue", meta.GetShortName())},
				"Skipping label '%s' as it already exists", apply)
		} else {
			meta.AddLabel(apply)
		}
	}
	remove := a.rule.Remove
	if remove != "" {
		if !meta.HasLabel(remove) {
			a.logger.DebugWith(
				log.MetaFields{log.F("issue", meta.GetShortName())},
				"Skipping label '%s' removal as it does not exists", remove)
		} else {
			meta.RemoveLabel(remove)
		}
	}
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) actions.Action {
	item := rule{}
	logger := log.Get("labeler")
	if e := mapstructure.Decode(config, &item); e != nil {
		logger.ErrorWith(log.MetaFields{log.E(e)}, "Unable to build action")
		return nil
	}
	return &action{rule: &item, logger: logger}
}

func init() {
	actions.RegisterAction("labeler", &factory{})
}
