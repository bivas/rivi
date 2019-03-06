package env

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bivas/rivi/pkg/types"
	"github.com/bivas/rivi/pkg/util/log"
)

type tempFSEnvironment struct {
	dir       string
	rulesFile string

	logger log.Logger
}

func (e *tempFSEnvironment) ConfigFilePath() string {
	return e.rulesFile
}

func (e *tempFSEnvironment) Create(data types.ReadOnlyData) error {
	temp, err := ioutil.TempDir("", "rivi-env-")
	if err != nil {
		return err
	}
	e.logger.DebugWith(log.MetaFields{
		log.F("issue", data.GetShortName()),
		log.F("dir", temp)}, "Created temp dir")
	e.dir = temp
	rules := filepath.Join(e.dir, types.RulesConfigFileName)
	if err := ioutil.WriteFile(
		rules,
		[]byte(data.GetRepository().GetRulesFile()),
		0600); err != nil {
		return err
	}
	e.rulesFile = rules
	return nil
}

func (e *tempFSEnvironment) Cleanup() error {
	e.logger.DebugWith(
		log.MetaFields{
			log.F("dir", e.dir)}, "Cleanup temp dir")
	return os.RemoveAll(e.dir)
}

func tempFSEnvironmentProvider() Environment {
	return &tempFSEnvironment{logger: le.Get("tempfs")}
}
