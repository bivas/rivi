package bot

import (
	"errors"
	"github.com/bivas/rivi/config"
	"github.com/bivas/rivi/engine"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/bivas/rivi/util/state"
)

type JobHandler interface {
	Handle(<-chan types.Data)
}

type localJobHandler struct {
	logger log.Logger
}

func (h *localJobHandler) internalHandle(data types.Data) error {
	env, err := GetEnvironment()
	if err != nil {
		h.logger.ErrorWith(
			log.MetaFields{log.E(err), log.F("issue", data.GetShortName())},
			"Failed to get environment")
		return err
	}
	defer env.Cleanup()
	if err := env.Create(data); err != nil {
		h.logger.ErrorWith(
			log.MetaFields{log.E(err), log.F("issue", data.GetShortName())},
			"Failed to create environment")
		return err
	}
	c, err := config.NewConfiguration(env.ConfigFilePath())
	if err != nil {
		h.logger.ErrorWith(
			log.MetaFields{log.E(err), log.F("issue", data.GetShortName())},
			"Failed to create configuration")
		return err
	}
	meta, ok := types.BuildComplete(c.GetClientConfig(), data)
	if !ok {
		return errors.New("Nothing to process")
	}
	applied := engine.ProcessRules(c.GetRules(), state.New(c, meta))
	h.logger.DebugWith(
		log.MetaFields{log.F("rules", applied)}, "Applied rules")
	return nil
}

func (h *localJobHandler) Handle(incoming <-chan types.Data) {
	for {
		data, ok := <-incoming
		if !ok {
			h.logger.Info("Stopping job handler")
			break
		}
		h.logger.InfoWith(log.MetaFields{log.F("data", data.GetShortName())}, "Got data from job channel")
		h.internalHandle(data)
	}
}
