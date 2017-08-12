package bot

import (
	"errors"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"gopkg.in/src-d/go-billy.v3/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	RULES_CONFIG_FILE = ".rivi.rules.yaml"
)

var (
	le = log.Get("env")
)

type Environment interface {
	Create(types.Data) error
	Cleanup() error
	ConfigFilePath() string
}

type tempFSEnvironment struct {
	dir       string
	rulesFile string

	logger log.Logger
}

func (e *tempFSEnvironment) ConfigFilePath() string {
	return e.rulesFile
}

func (e *tempFSEnvironment) clone(data types.Data) error {
	fs := osfs.New(e.dir)
	p, _ := filesystem.NewStorage(fs)
	repo, err := git.Clone(p, fs, &git.CloneOptions{
		URL:               data.GetOrigin().GitURL,
		ReferenceName:     plumbing.ReferenceName("refs/heads/" + data.GetOrigin().Ref),
		SingleBranch:      true,
		Depth:             1,
		RecurseSubmodules: git.NoRecurseSubmodules,
		Progress:          nil,
	})
	if err != nil {
		return err
	}
	head, err := repo.Head()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(head.Hash().String(), data.GetOrigin().Head) {
		return errors.New("Head Ref and Origin Ref do not match")
	}
	return err
}

func (e *tempFSEnvironment) Create(data types.Data) error {
	temp, err := ioutil.TempDir("", "rivi-env-")
	if err != nil {
		return err
	}
	e.logger.DebugWith(log.MetaFields{log.F("dir", temp)}, "Created temp dir")
	e.dir = temp
	if err := e.clone(data); err != nil {
		return err
	}
	rules := filepath.Join(e.dir, RULES_CONFIG_FILE)
	if _, err := os.Stat(rules); err != nil {
		return err
	}
	e.rulesFile = rules
	return nil
}

func (e *tempFSEnvironment) Cleanup() error {
	return os.RemoveAll(e.dir)
}

func tempFSEnvironmentProvider() Environment {
	return &tempFSEnvironment{logger: le.Get("temp")}
}

type EnvironmentProvider func() HookListenerQueue

var defaultEnvironmentProvider = tempFSEnvironmentProvider

func GetEnvironment() (Environment, error) {
	return defaultEnvironmentProvider(), nil
}
