package status

import (
	"errors"

	"github.com/bivas/rivi/pkg/engine/actions"
	"github.com/bivas/rivi/pkg/types"
	"github.com/bivas/rivi/pkg/util/log"

	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/multistep"
	"github.com/prometheus/client_golang/prometheus"
)

type action struct {
	rule *rule
	err  error

	logger log.Logger
}

type HasSetStatusAPIData interface {
	SetStatus(string, types.State)
}

func (a *action) Apply(state multistep.StateBag) {
	meta := state.Get("data").(types.Data)
	setStatus, ok := meta.(HasSetStatusAPIData)
	if !ok {
		a.logger.Warning("Event data does not support setting status. Check your configurations")
		a.err = errors.New("event data does not support setting status")
		return
	}
	setStatus.SetStatus(a.rule.Description, types.GetState(a.rule.State))
}

type factory struct {
}

func (*factory) BuildAction(config map[string]interface{}) actions.Action {
	item := rule{}
	logger := log.Get("status")
	if e := mapstructure.Decode(config, &item); e != nil {
		logger.ErrorWith(log.MetaFields{log.E(e)}, "Unable to build action")
		return nil
	}
	item.Defaults()
	return &action{rule: &item, logger: logger}
}

var counter = actions.NewCounter("status")

func init() {
	actions.RegisterAction("status", &factory{})
	prometheus.Register(counter)
}
