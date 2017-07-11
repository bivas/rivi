package bot

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"sort"

	"github.com/bivas/rivi/util"
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
		util.Logger.Warning("Request for namespace '%s' matched nothing", namespace)
		return nil, fmt.Errorf("Request for namespace '%s' matched nothing", namespace)
	}
	return configuration, nil
}

func (b *bot) getIssueLock(namespaceLock *sync.Mutex, data EventData) *sync.Mutex {
	defer namespaceLock.Unlock()
	id := fmt.Sprintf("%s/%s#%d", data.GetOwner(), data.GetRepo(), data.GetNumber())
	util.Logger.Debug("acquire namespace lock during rules process")
	issueLocker, exists := b.repoIssueMutexes.Get(id)
	if !exists {
		issueLocker = &sync.Mutex{}
		b.repoIssueMutexes.Set(id, issueLocker, cache.DefaultExpiration)
	}
	util.Logger.Debug("acquire repo issue %s lock during rules process", id)
	issueLocker.(*sync.Mutex).Lock()
	return issueLocker.(*sync.Mutex)
}

type rulesGroup struct {
	key   int
	rules []Rule
}

type rulesGroups []rulesGroup

func (r rulesGroups) Len() int {
	return len(r)
}

func (r rulesGroups) Less(i, j int) bool {
	return r[i].key < r[j].key
}

func (r rulesGroups) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func groupByRuleOrder(rules []Rule) []rulesGroup {
	groupIndexes := make(map[int]rulesGroup)
	for _, rule := range rules {
		key := rule.Order()
		rules, exists := groupIndexes[key]
		if !exists {
			rules = rulesGroup{key, make([]Rule, 0)}
		}
		rules.rules = append(rules.rules, rule)
		groupIndexes[key] = rules
	}
	util.Logger.Debug("%d Rules are grouped to %d rule groups", len(rules), len(groupIndexes))
	groupResult := make([]rulesGroup, 0)
	for _, group := range groupIndexes {
		groupResult = append(groupResult, group)
	}
	sort.Sort(rulesGroups(groupResult))
	return groupResult
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
		util.Logger.Debug("Skipping rule processing for %d (couldn't build complete data)", partial.GetNumber())
		return result
	}
	for _, group := range rules {
		util.Logger.Debug("Processing rule group of %d order with %d rules", group.key, len(group.rules))
		for _, rule := range group.rules {
			if rule.Accept(data) {
				util.Logger.Debug("Accepting rule %s for '%s'", rule.Name(), data.GetTitle())
				applied = append(applied, rule)
				result.AppliedRules = append(result.AppliedRules, rule.Name())
			}
		}
		for _, rule := range applied {
			util.Logger.Debug("Applying rule %s for '%s'", rule.Name(), data.GetTitle())
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
	util.Logger.Debug("acquire global lock during namespace process")
	locker, exists := b.namespaceMutexes.Get(namespace)
	if !exists {
		locker = &sync.Mutex{}
		b.namespaceMutexes.Set(namespace, locker, cache.DefaultExpiration)
	}
	util.Logger.Debug("acquire namespace '%s' lock", namespace)
	locker.(*sync.Mutex).Lock()
	util.Logger.Debug("release global lock during namespace process")
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
	util.Logger.Debug("release namespace %s lock", namespace)
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
