package runner

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bivas/rivi/pkg/config"
	"github.com/bivas/rivi/pkg/engine"
	"github.com/bivas/rivi/pkg/types"
	"github.com/bivas/rivi/pkg/types/builder"
	"github.com/bivas/rivi/pkg/util/log"
	"github.com/bivas/rivi/pkg/util/state"
	"github.com/patrickmn/go-cache"
)

var runnerLog = log.Get("runner")

type HandledEventResult struct {
	AppliedRules []string `json:"applied_rules,omitempty"`
	Message      string   `json:"message,omitempty"`
}

type Runner interface {
	HandleEvent(r *http.Request) *HandledEventResult
}

type runnable struct {
	defaultNamespace string
	configurations   map[string]config.Configuration

	globalLocker     *sync.Mutex
	namespaceMutexes *cache.Cache
	repoIssueMutexes *cache.Cache
}

func (b *runnable) getCurrentConfiguration(namespace string) (config.Configuration, error) {
	if namespace == "" {
		namespace = b.defaultNamespace
	}
	configuration, exists := b.configurations[namespace]
	if !exists {
		runnerLog.Warning("Request for namespace '%s' matched nothing", namespace)
		return nil, fmt.Errorf("Request for namespace '%s' matched nothing", namespace)
	}
	return configuration, nil
}

func (b *runnable) getIssueLock(namespaceLock *sync.Mutex, data types.InfoData) *sync.Mutex {
	defer namespaceLock.Unlock()
	id := data.GetShortName()
	runnerLog.Debug("acquire namespace lock during rules process")
	issueLocker, exists := b.repoIssueMutexes.Get(id)
	if !exists {
		issueLocker = &sync.Mutex{}
		b.repoIssueMutexes.Set(id, issueLocker, cache.DefaultExpiration)
	}
	runnerLog.Debug("acquire repo issue %s lock during rules process", id)
	issueLocker.(*sync.Mutex).Lock()
	return issueLocker.(*sync.Mutex)
}

func (b *runnable) processRules(namespaceLock *sync.Mutex, config config.Configuration, partial types.HookData) *HandledEventResult {
	issueLocker := b.getIssueLock(namespaceLock, partial)
	defer issueLocker.Unlock()

	meta, ok := builder.BuildComplete(config.GetClientConfig(), partial)
	if !ok {
		runnerLog.Debug("Skipping rule processing for %s (couldn't build complete data)", partial.GetShortName())
		return &HandledEventResult{
			AppliedRules: []string{},
		}
	}
	return &HandledEventResult{
		AppliedRules: engine.ProcessRules(config.GetRules(), state.New(config, meta)),
	}
}

func (b *runnable) HandleEvent(r *http.Request) *HandledEventResult {
	namespace := r.URL.Query().Get("namespace")
	b.globalLocker.Lock()
	runnerLog.Debug("acquire global lock during namespace process")
	locker, exists := b.namespaceMutexes.Get(namespace)
	if !exists {
		locker = &sync.Mutex{}
		b.namespaceMutexes.Set(namespace, locker, cache.DefaultExpiration)
	}
	runnerLog.Debug("acquire namespace '%s' lock", namespace)
	locker.(*sync.Mutex).Lock()
	runnerLog.Debug("release global lock during namespace process")
	b.globalLocker.Unlock()
	workingConfiguration, err := b.getCurrentConfiguration(namespace)
	if err != nil {
		locker.(*sync.Mutex).Unlock()
		return &HandledEventResult{Message: err.Error()}
	}
	meta, process := builder.BuildFromHook(workingConfiguration.GetClientConfig(), r)
	if !process {
		locker.(*sync.Mutex).Unlock()
		return &HandledEventResult{Message: "Skipping rules processing (could be not supported event type)"}
	}
	runnerLog.Debug("release namespace '%s' lock", namespace)
	return b.processRules(locker.(*sync.Mutex), workingConfiguration, meta)
}

func New(configPaths ...string) (Runner, error) {
	b := &runnable{
		configurations:   make(map[string]config.Configuration),
		globalLocker:     &sync.Mutex{},
		namespaceMutexes: cache.New(time.Minute, 30*time.Second),
		repoIssueMutexes: cache.New(time.Minute, 20*time.Second),
	}
	for index, configPath := range configPaths {
		baseConfigPath := filepath.Base(configPath)
		namespace := strings.TrimSuffix(baseConfigPath, filepath.Ext(baseConfigPath))
		runnerLog.Debug("Loading configuration for namespace '%s'", namespace)
		if index == 0 {
			b.defaultNamespace = namespace
		}
		configuration, err := config.NewConfiguration(configPath)
		if err != nil {
			return nil, fmt.Errorf("Reading %s caused an error. %s", configPath, err)
		}
		b.configurations[namespace] = configuration
	}
	if len(b.configurations) == 0 {
		return nil, fmt.Errorf("Runner has no readable configuration!")
	}
	runnerLog.Debug("Runner is ready %+v", *b)
	return b, nil
}
