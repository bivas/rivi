package bot

import (
	"fmt"
	"github.com/bivas/rivi/util"
	"github.com/patrickmn/go-cache"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type HandledEventResult struct {
	AppliedRules []string `json:"applied_rules,omitempty"`
	Message      string   `json:"message,omitempty"`
}

type Bot interface {
	HandleEvent(r *http.Request) *HandledEventResult
}

type bot struct {
	defaultNamespace string
	configurations   map[string]Configuration

	cacheLocker      *sync.Mutex
	namespaceMutexes *cache.Cache
	repoIssueMutexes *cache.Cache
}

func (b *bot) getCurrentConfiguration(namespace string) (Configuration, error) {
	if namespace == "" {
		namespace = b.defaultNamespace
	}
	configuration, exists := b.configurations[namespace]
	if !exists {
		util.Logger.Warning("Request for namespace '%s' matched nothing", namespace)
		return nil, fmt.Errorf("Request for namespace '%s' matched nothing", namespace)
	}
	return configuration, nil
}

func (b *bot) processRules(configuration Configuration, data EventData) *HandledEventResult {
	id := fmt.Sprintf("%s/%s#%d", data.GetOwner(), data.GetRepo(), data.GetNumber())
	util.Logger.Debug("acquire global lock during rules process")
	b.cacheLocker.Lock()
	locker, exists := b.repoIssueMutexes.Get(id)
	if !exists {
		locker = &sync.Mutex{}
		b.repoIssueMutexes.Set(id, locker, cache.DefaultExpiration)
	}
	util.Logger.Debug("acquire repo issue %s lock during rules process", id)
	locker.(*sync.Mutex).Lock()
	defer locker.(*sync.Mutex).Unlock()
	util.Logger.Debug("release global lock during rules process")
	b.cacheLocker.Unlock()
	applied := make([]Rule, 0)
	result := &HandledEventResult{
		AppliedRules: []string{},
	}
	for _, rule := range configuration.GetRules() {
		if rule.Accept(data) {
			util.Logger.Debug("Accepting rule %s for '%s'", rule.Name(), data.GetTitle())
			applied = append(applied, rule)
			result.AppliedRules = append(result.AppliedRules, rule.Name())
		}
	}
	for _, rule := range applied {
		util.Logger.Debug("Applying rule %s for '%s'", rule.Name(), data.GetTitle())
		for _, action := range rule.Actions() {
			action.Apply(configuration, data)
		}
	}
	return result
}

func (b *bot) HandleEvent(r *http.Request) *HandledEventResult {
	namespace := r.URL.Query().Get("namespace")
	b.cacheLocker.Lock()
	util.Logger.Debug("acquire global lock during namespace process")
	locker, exists := b.namespaceMutexes.Get(namespace)
	if !exists {
		locker = &sync.Mutex{}
		b.namespaceMutexes.Set(namespace, locker, cache.DefaultExpiration)
	}
	util.Logger.Debug("acquire namespace %s lock", namespace)
	locker.(*sync.Mutex).Lock()
	util.Logger.Debug("release global lock during namespace process")
	b.cacheLocker.Unlock()
	workingConfiguration, err := b.getCurrentConfiguration(namespace)
	if err != nil {
		return &HandledEventResult{Message: err.Error()}
	}
	data, process := buildFromRequest(workingConfiguration.GetClientConfig(), r)
	if !process {
		return &HandledEventResult{Message: "Skipping rules processing (could be not supported event type)"}
	}
	util.Logger.Debug("release namespace %s lock", namespace)
	locker.(*sync.Mutex).Unlock()
	return b.processRules(workingConfiguration, data)
}

func New(configPaths ...string) (Bot, error) {
	b := &bot{
		configurations:   make(map[string]Configuration),
		cacheLocker:      &sync.Mutex{},
		namespaceMutexes: cache.New(time.Minute, 30*time.Second),
		repoIssueMutexes: cache.New(time.Minute, 20*time.Second),
	}
	for index, configPath := range configPaths {
		baseConfigPath := filepath.Base(configPath)
		namespace := strings.TrimSuffix(baseConfigPath, filepath.Ext(baseConfigPath))
		util.Logger.Debug("Loading configuration for namespace '%s'", namespace)
		if index == 0 {
			b.defaultNamespace = namespace
		}
		configuration, err := newConfiguration(configPath)
		if err != nil {
			return nil, fmt.Errorf("Reading %s caused an error. %s", configPath, err)
		}
		b.configurations[namespace] = configuration
	}
	if len(b.configurations) == 0 {
		return nil, fmt.Errorf("Bot has no readable configuration!")
	}
	util.Logger.Debug("Bot is ready %+v", *b)
	return b, nil
}
