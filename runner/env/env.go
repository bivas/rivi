package env

import (
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
)

var (
	le = log.Get("env")
)

type Environment interface {
	Create(types.ReadOnlyData) error
	Cleanup() error
	ConfigFilePath() string
}

type EnvironmentProvider func() Environment

var defaultEnvironmentProvider EnvironmentProvider = tempFSEnvironmentProvider

func GetEnvironment() (Environment, error) {
	return defaultEnvironmentProvider(), nil
}
