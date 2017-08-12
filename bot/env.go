package bot

import (
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"gopkg.in/src-d/go-billy.v3/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	RULES_CONFIG_FILE = ".rivi.rules.yaml"
)

type Environment interface {
	Create(types.Data) (string, error)
	Cleanup() error
}

type tempFSEnvironment struct {
	path string
}

func (e *tempFSEnvironment) clone(data types.Data) error {
	fs := osfs.New(e.path)
	p, _ := filesystem.NewStorage(fs)
	_, err := git.Clone(p, fs, &git.CloneOptions{
		URL: data.GetOrigin().GitURL,
		ReferenceName: plumbing.NewHashReference(
			plumbing.ReferenceName(plumbing.HEAD),
			plumbing.NewHash(data.GetOrigin().Ref)).Name(),
		SingleBranch:      true,
		Depth:             1,
		RecurseSubmodules: git.NoRecurseSubmodules,
		Progress:          nil,
	})
	return err
}

func (e *tempFSEnvironment) Create(data types.Data) (string, error) {
	temp, err := ioutil.TempDir("", "rivi-env-")
	if err != nil {
		return "", err
	}
	log.DebugWith(log.MetaFields{log.F("path", temp)}, "Created temp path")
	e.path = temp
	if e.clone(data) != nil {
		return "", err
	}
	rules := filepath.Join(e.path, RULES_CONFIG_FILE)
	if _, err := os.Stat(rules); err != nil {
		return "", err
	}
	return rules, nil
}

func (e *tempFSEnvironment) Cleanup() error {
	return os.RemoveAll(e.path)
}

func tempFSEnvironmentProvider() Environment {
	return &tempFSEnvironment{}
}

type EnvironmentProvider func() HookListenerQueue

var defaultEnvironmentProvider = tempFSEnvironmentProvider

func GetEnvironment() (Environment, error) {
	return defaultEnvironmentProvider(), nil
}
