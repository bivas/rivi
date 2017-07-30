package bot

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bivas/rivi/util/log"
	"github.com/patrickmn/go-cache"
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

	globalLocker     *sync.Mutex
	namespaceMutexes *cache.Cache
	repoIssueMutexes *cache.Cache
}

func (b *bot) getCurrentConfiguration(namespace string) (Configuration, error) {
	if namespace == "" {
		namespace = b.defaultNamespace
	}
	configuration, exists := b.configurations[namespace]
	if !exists {
		log.Warning("Request for namespace '%s' matched nothing", namespace)
		return nil, fmt.Errorf("Request for namespace '%s' matched nothing", namespace)
	}
	return configuration, nil
}

func (b *bot) getIssueLock(namespaceLock *sync.Mutex, data EventData) *sync.Mutex {
	defer namespaceLock.Unlock()
	id := data.GetShortName()
	log.Debug("acquire namespace lock during rules process")
	issueLocker, exists := b.repoIssueMutexes.Get(id)
	if !exists {
		issueLocker = &sync.Mutex{}
		b.repoIssueMutexes.Set(id, issueLocker, cache.DefaultExpiration)
	}
	log.Debug("acquire repo issue %s lock during rules process", id)
	issueLocker.(*sync.Mutex).Lock()
	return issueLocker.(*sync.Mutex)
}

func (b *bot) processRules(namespaceLock *sync.Mutex, config Configuration, partial EventData, r *http.Request) *HandledEventResult {
	rules := groupByRuleOrder(config.GetRules())
	issueLocker := b.getIssueLock(namespaceLock, partial)
	defer issueLocker.Unlock()

	applied := make([]Rule, 0)
	result := &HandledEventResult{
		AppliedRules: []string{},
	}
	data, ok := completeBuild(config.GetClientConfig(), r, partial)
	if !ok {
		log.Debug("Skipping rule processing for %s (couldn't build complete data)", partial.GetShortName())
		return result
	}
	for _, group := range rules {
		log.DebugWith(log.MetaFields{log.F("issue", data.GetShortName())}, "Processing rule group of %d order with %d rules", group.key, len(group.rules))
		for _, rule := range group.rules {
			if rule.Accept(data) {
				log.DebugWith(log.MetaFields{log.F("issue", data.GetShortName())}, "Accepting rule %s", rule.Name())
				applied = append(applied, rule)
				result.AppliedRules = append(result.AppliedRules, rule.Name())
			}
		}
		for _, rule := range applied {
			log.DebugWith(log.MetaFields{log.F("issue", data.GetShortName())}, "Applying rule %s", rule.Name())
			for _, action := range rule.Actions() {
				action.Apply(config, data)
			}
		}
	}
	return result
}

func (b *bot) HandleEvent(r *http.Request) *HandledEventResult {
	namespace := r.URL.Query().Get("namespace")
	b.globalLocker.Lock()
	log.Debug("acquire global lock during namespace process")
	locker, exists := b.namespaceMutexes.Get(namespace)
	if !exists {
		locker = &sync.Mutex{}
		b.namespaceMutexes.Set(namespace, locker, cache.DefaultExpiration)
	}
	log.Debug("acquire namespace '%s' lock", namespace)
	locker.(*sync.Mutex).Lock()
	log.Debug("release global lock during namespace process")
	b.globalLocker.Unlock()
	workingConfiguration, err := b.getCurrentConfiguration(namespace)
	if err != nil {
		locker.(*sync.Mutex).Unlock()
		return &HandledEventResult{Message: err.Error()}
	}
	data, process := buildFromRequest(workingConfiguration.GetClientConfig(), r)
	if !process {
		locker.(*sync.Mutex).Unlock()
		return &HandledEventResult{Message: "Skipping rules processing (could be not supported event type)"}
	}
	log.Debug("release namespace '%s' lock", namespace)
	return b.processRules(locker.(*sync.Mutex), workingConfiguration, data, r)
}

func New(configPaths ...string) (Bot, error) {
	b := &bot{
		configurations:   make(map[string]Configuration),
		globalLocker:     &sync.Mutex{},
		namespaceMutexes: cache.New(time.Minute, 30*time.Second),
		repoIssueMutexes: cache.New(time.Minute, 20*time.Second),
	}
	for index, configPath := range configPaths {
		baseConfigPath := filepath.Base(configPath)
		namespace := strings.TrimSuffix(baseConfigPath, filepath.Ext(baseConfigPath))
		log.Debug("Loading configuration for namespace '%s'", namespace)
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
	log.Debug("Bot is ready %+v", *b)
	return b, nil
}
