package bot

import (
	"github.com/bivas/rivi/types"
	"io/ioutil"
)

type Environment interface {
	Create() (string, error)
	Cleanup() error
}

type tempFSEnvironment struct {
	path string
}

func (e *tempFSEnvironment) Create() (string, error) {
	temp, err := ioutil.TempDir("", "rivi-env-")
	if err != nil {
		return "", err
	}
	return temp, nil
}

func (e *tempFSEnvironment) Cleanup() error {
	panic("implement me")
}

func tempFSEnvironmentProvider() Environment {
	return &tempFSEnvironment{}
}

type EnvironmentProvider func() HookListenerQueue

var defaultEnvironementProvider = tempFSEnvironmentProvider

func GetEnvironment(data types.Data) (Environment, error) {
	return defaultEnvironementProvider(), nil
}
