package local

import (
	"errors"

	"github.com/bivas/rivi/pkg/config"
	"github.com/bivas/rivi/pkg/engine"
	"github.com/bivas/rivi/pkg/runner/env"
	"github.com/bivas/rivi/pkg/runner/types"
	"github.com/bivas/rivi/pkg/types/builder"
	"github.com/bivas/rivi/pkg/util/log"
	"github.com/bivas/rivi/pkg/util/state"
	"github.com/prometheus/client_golang/prometheus"
)

type workUnit struct {
	incoming <-chan *types.Message
	logger   log.Logger
}

func (w *workUnit) internalHandle(msg *types.Message) error {
	environment, err := env.GetEnvironment()
	if err != nil {
		w.logger.ErrorWith(
			log.MetaFields{log.E(err), log.F("issue", msg.Data.GetShortName())},
			"Failed to get environment")
		return err
	}
	defer environment.Cleanup()
	meta, ok := builder.BuildComplete(msg.Config, msg.Data)
	if !ok {
		return errors.New("nothing to process")
	}
	if err := environment.Create(meta); err != nil {
		w.logger.ErrorWith(
			log.MetaFields{log.E(err), log.F("issue", msg.Data.GetShortName())},
			"Failed to create environment")
		return err
	}
	c, err := config.NewConfiguration(environment.ConfigFilePath())
	if err != nil {
		w.logger.ErrorWith(
			log.MetaFields{log.E(err), log.F("issue", msg.Data.GetShortName())},
			"Failed to create configuration")
		return err
	}
	applied := engine.ProcessRules(c.GetRules(), state.New(c, meta))
	w.logger.DebugWith(
		log.MetaFields{log.F("rules", applied)}, "Applied rules")
	return nil
}

func (w *workUnit) Handle() {
	for {
		msg, ok := <-w.incoming
		if !ok {
			w.logger.Info("Stopping job handler")
			break
		}
		w.logger.InfoWith(
			log.MetaFields{
				log.F("data", msg.Data.GetShortName()),
			}, "Got data from job channel")
		timer := prometheus.NewTimer(handleHistogram)
		if err := w.internalHandle(msg); err != nil {
			handleErrorCounter.Inc()
			w.logger.WarningWith(log.MetaFields{
				log.E(err),
				log.F("data", msg.Data.GetShortName())}, "Error when handling data")
		}
		timer.ObserveDuration()
	}
}

var (
	handleHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "rivi",
		Subsystem: "workunit",
		Name:      "handle",
		Help:      "Measure handling of event data",
	})
	handleErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "rivi",
		Subsystem: "workunit",
		Name:      "failure",
		Help:      "Failure to handle event data",
	})
)

func init() {
	prometheus.Register(handleHistogram)
	prometheus.Register(handleErrorCounter)
}
