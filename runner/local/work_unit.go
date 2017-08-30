package local

import (
	"errors"

	"github.com/bivas/rivi/config"
	"github.com/bivas/rivi/engine"
	"github.com/bivas/rivi/runner/env"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/types/builder"
	"github.com/bivas/rivi/util/log"
	"github.com/bivas/rivi/util/state"
)

type workUnit struct {
	incoming <-chan types.HookData
	logger   log.Logger
}

func (w *workUnit) internalHandle(data types.HookData) error {
	environment, err := env.GetEnvironment()
	if err != nil {
		w.logger.ErrorWith(
			log.MetaFields{log.E(err), log.F("issue", data.GetShortName())},
			"Failed to get environment")
		return err
	}
	defer environment.Cleanup()
	c, err := config.NewConfiguration(environment.ConfigFilePath())
	if err != nil {
		w.logger.ErrorWith(
			log.MetaFields{log.E(err), log.F("issue", data.GetShortName())},
			"Failed to create configuration")
		return err
	}
	meta, ok := builder.BuildComplete(c.GetClientConfig(), data)
	if !ok {
		return errors.New("Nothing to process")
	}
	if err := environment.Create(meta); err != nil {
		w.logger.ErrorWith(
			log.MetaFields{log.E(err), log.F("issue", data.GetShortName())},
			"Failed to create environment")
		return err
	}
	applied := engine.ProcessRules(c.GetRules(), state.New(c, meta))
	w.logger.DebugWith(
		log.MetaFields{log.F("rules", applied)}, "Applied rules")
	return nil
}

func (w *workUnit) Handle() {
	for {
		data, ok := <-w.incoming
		if !ok {
			w.logger.Info("Stopping job handler")
			break
		}
		w.logger.InfoWith(log.MetaFields{log.F("data", data.GetShortName())}, "Got data from job channel")
		if err := w.internalHandle(data); err != nil {
			w.logger.WarningWith(log.MetaFields{
				log.E(err),
				log.F("data", data.GetShortName())}, "Error when handling data")
		}
	}
}
